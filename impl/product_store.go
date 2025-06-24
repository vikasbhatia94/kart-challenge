package impl

import "backend-challenge/api"

// productStore simulates an in-memory product database
var productStore = []*api.Product{
	{
		ID:       ptrString("1"),
		Name:     ptrString("Waffle with Berries"),
		Category: ptrString("Waffle"),
		Price:    ptrFloat(6.5),
		Image: &api.ImageVariants{
			Thumbnail: ptrString("https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg"),
			Mobile:    ptrString("https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg"),
			Tablet:    ptrString("https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg"),
			Desktop:   ptrString("https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg"),
		},
	},
	{
		ID:       ptrString("2"),
		Name:     ptrString("Vanilla Bean Crème Brûlée"),
		Category: ptrString("Crème Brûlée"),
		Price:    ptrFloat(7),
		Image: &api.ImageVariants{
			Thumbnail: ptrString("https://orderfoodonline.deno.dev/public/images/image-creme-brulee-thumbnail.jpg"),
			Mobile:    ptrString("https://orderfoodonline.deno.dev/public/images/image-creme-brulee-mobile.jpg"),
			Tablet:    ptrString("https://orderfoodonline.deno.dev/public/images/image-creme-brulee-tablet.jpg"),
			Desktop:   ptrString("https://orderfoodonline.deno.dev/public/images/image-creme-brulee-desktop.jpg"),
		},
	},
	{
		ID:       ptrString("3"),
		Name:     ptrString("Macaron Mix of Five"),
		Category: ptrString("Macaron"),
		Price:    ptrFloat(8),
		Image: &api.ImageVariants{
			Thumbnail: ptrString("https://orderfoodonline.deno.dev/public/images/image-macaron-thumbnail.jpg"),
			Mobile:    ptrString("https://orderfoodonline.deno.dev/public/images/image-macaron-mobile.jpg"),
			Tablet:    ptrString("https://orderfoodonline.deno.dev/public/images/image-macaron-tablet.jpg"),
			Desktop:   ptrString("https://orderfoodonline.deno.dev/public/images/image-macaron-desktop.jpg"),
		},
	},
	{
		ID:       ptrString("4"),
		Name:     ptrString("Classic Tiramisu"),
		Category: ptrString("Tiramisu"),
		Price:    ptrFloat(5.5),
		Image: &api.ImageVariants{
			Thumbnail: ptrString("https://orderfoodonline.deno.dev/public/images/image-tiramisu-thumbnail.jpg"),
			Mobile:    ptrString("https://orderfoodonline.deno.dev/public/images/image-tiramisu-mobile.jpg"),
			Tablet:    ptrString("https://orderfoodonline.deno.dev/public/images/image-tiramisu-tablet.jpg"),
			Desktop:   ptrString("https://orderfoodonline.deno.dev/public/images/image-tiramisu-desktop.jpg"),
		},
	},
	{
		ID:       ptrString("5"),
		Name:     ptrString("Pistachio Baklava"),
		Category: ptrString("Baklava"),
		Price:    ptrFloat(4),
		Image: &api.ImageVariants{
			Thumbnail: ptrString("https://orderfoodonline.deno.dev/public/images/image-baklava-thumbnail.jpg"),
			Mobile:    ptrString("https://orderfoodonline.deno.dev/public/images/image-baklava-mobile.jpg"),
			Tablet:    ptrString("https://orderfoodonline.deno.dev/public/images/image-baklava-tablet.jpg"),
			Desktop:   ptrString("https://orderfoodonline.deno.dev/public/images/image-baklava-desktop.jpg"),
		},
	},
}

func ListAllProducts() []*api.Product {
	return productStore
}

func GetProductByID(id string) *api.Product {
	for _, p := range productStore {
		if p.ID != nil && *p.ID == id {
			return p
		}
	}
	return nil
}
