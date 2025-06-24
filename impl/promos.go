package impl

import (
	"bufio"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ValidPromoCodes = make(map[string]struct{})
	promoOnce       sync.Once
	ForceRefresh    = flag.Bool("refresh-promos", false, "Force refresh promo code cache")
)

const cacheFile = "valid_promos.gob"

func DownloadPromoFiles(promoPath string, files []string) error {
	urls := map[string]string{
		"couponbase1.gz": "https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase1.gz",
		"couponbase2.gz": "https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase2.gz",
		"couponbase3.gz": "https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase3.gz",
	}

	for _, file := range files {
		url, ok := urls[file]
		if !ok {
			log.Printf("No URL found for %s\n", file)
			continue
		}
		filePath := filepath.Join(promoPath, file)
		log.Printf("Downloading %s...\n", file)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to download %s: %v\n", file, err)
			continue
		}
		defer resp.Body.Close()

		out, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file %s: %v\n", filePath, err)
			continue
		}
		if _, err := io.Copy(out, resp.Body); err != nil {
			log.Printf("Failed to save %s: %v\n", file, err)
		}
		out.Close()
		log.Printf("Downloaded %s\n", file)
	}
	return nil
}

func LoadPromoCodes(folderPath string) error {
	// If cache already exists, skip processing
	if !*ForceRefresh && fileExists(cacheFile) {
		if loadFromCache(cacheFile) == nil {
			return nil
		}
	}
	var loadErr error
	promoOnce.Do(func() {
		start := time.Now()

		if !*ForceRefresh && loadFromCache(cacheFile) == nil {
			return
		}

		log.Println("Processing promo code files (preferring .txt over .gz)...")

		bases := []string{"couponbase1", "couponbase2", "couponbase3"}
		tempFiles := make([]string, 0, 3)
		var dedupWg sync.WaitGroup

		// Determine the number of concurrent workers. Use the number of files
		// or the number of CPU cores, whichever is smaller. This avoids creating
		// more goroutines than necessary.
		numFiles := len(bases)
		workerCount := runtime.NumCPU()
		if numFiles < workerCount {
			workerCount = numFiles
		}

		sem := make(chan struct{}, workerCount)
		mu := sync.Mutex{}

		for _, base := range bases {
			txtExists := fileExists(filepath.Join(folderPath, base+".txt"))
			gzExists := fileExists(filepath.Join(folderPath, base+".gz"))
			if !txtExists && !gzExists {
				log.Println("Missing promo file:", base)
				DownloadPromoFiles(folderPath, []string{base + ".gz"})
			}
			var path string
			if fileExists(filepath.Join(folderPath, base+".txt")) {
				path = filepath.Join(folderPath, base+".txt")
			} else if fileExists(filepath.Join(folderPath, base+".gz")) {
				path = filepath.Join(folderPath, base+".gz")
			} else {
				loadErr = fmt.Errorf("missing promo file for base: %s", base)
				return
			}

			rawPath := filepath.Join(folderPath, base+".raw.tmp")
			sortedPath := filepath.Join(folderPath, base+".sorted.tmp")

			dedupWg.Add(1)
			sem <- struct{}{}
			go func(base, path, rawPath, sortedPath string) {
				defer dedupWg.Done()
				defer func() {
					<-sem
				}()

				f, err := os.Open(path)
				if err != nil {
					log.Printf("Failed to open %s: %v", path, err)
					return
				}
				scanner := bufio.NewScanner(f)
				rawFile, err := os.Create(rawPath)
				if err != nil {
					log.Printf("Failed to create raw temp file: %v", err)
					f.Close()
					return
				}
				for scanner.Scan() {
					words := strings.Fields(scanner.Text())
					for _, w := range words {
						w = strings.TrimSpace(w)
						if isAscii(w) && len(w) >= 8 && len(w) <= 10 {
							fmt.Fprintln(rawFile, w)
						}
					}
				}
				rawFile.Close()
				f.Close()

				input, err := os.Open(rawPath)
				if err != nil {
					log.Printf("Failed to reopen raw file: %v", err)
					return
				}
				tokens := map[string]struct{}{}
				scanner = bufio.NewScanner(input)
				for scanner.Scan() {
					tokens[scanner.Text()] = struct{}{}
				}
				input.Close()

				out, err := os.Create(sortedPath)
				if err != nil {
					log.Printf("Failed to write sorted temp file: %v", err)
					return
				}
				keys := make([]string, 0, len(tokens))
				for k := range tokens {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintln(out, k)
				}
				out.Close()

				_ = os.Remove(path)
				_ = os.Remove(rawPath)
				log.Printf("Deleted %s and %s", filepath.Base(path), filepath.Base(rawPath))

				mu.Lock()
				tempFiles = append(tempFiles, sortedPath)
				mu.Unlock()
			}(base, path, rawPath, sortedPath)
		}

		dedupWg.Wait()

		if len(tempFiles) < 3 {
			loadErr = errors.New("not enough valid temp promo files")
			return
		}

		// 3-way merge
		files := make([]*os.File, 3)
		scanners := make([]*bufio.Scanner, 3)
		current := make([]string, 3)
		active := make([]bool, 3)
		for i, file := range tempFiles[:3] {
			f, err := os.Open(file)
			if err != nil {
				loadErr = err
				return
			}
			files[i] = f
			scanners[i] = bufio.NewScanner(f)
			active[i] = scanners[i].Scan()
			if active[i] {
				current[i] = scanners[i].Text()
			}
		}
		for active[0] || active[1] || active[2] {
			min := ""
			for i := 0; i < 3; i++ {
				if active[i] && (min == "" || current[i] < min) {
					min = current[i]
				}
			}
			count := 0
			for i := 0; i < 3; i++ {
				if active[i] && current[i] == min {
					count++
					active[i] = scanners[i].Scan()
					if active[i] {
						current[i] = scanners[i].Text()
					}
				}
			}
			if count >= 2 {
				ValidPromoCodes[min] = struct{}{}
			}
		}
		for _, f := range files {
			f.Close()
		}
		for _, path := range tempFiles {
			_ = os.Remove(path)
		}
		_ = saveToCache(cacheFile)
		log.Printf("Loaded %d valid promo codes in %.2fs.\n", len(ValidPromoCodes), time.Since(start).Seconds())
	})
	return loadErr
}

func IsPromoCodeValid(code string) bool {
	_, ok := ValidPromoCodes[strings.TrimSpace(code)]
	return ok
}

func saveToCache(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(ValidPromoCodes)
}

func loadFromCache(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(&ValidPromoCodes)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isAscii(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}
