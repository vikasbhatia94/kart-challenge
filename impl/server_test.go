// impl/server_test.go
package impl

import (
	"backend-challenge/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestServer() *Server {
	ValidPromoCodes = map[string]struct{}{
		"HAPPYHRS": {},
		"FIFTYOFF": {},
	}
	return NewServer()
}

func TestListProducts(t *testing.T) {
	ts := setupTestServer()
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product", nil)
	ts.ListProducts(r, req)

	if r.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", r.Code)
	}

	var products []api.Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to parse response body: %v. Body: %s", err, r.Body.String())
	}
	if len(products) == 0 {
		t.Error("Expected at least one product, got 0")
	}
}

func TestGetProductValid(t *testing.T) {
	ts := setupTestServer()
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/1", nil)
	ts.GetProduct(r, req, 1)

	if r.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", r.Code)
	}
}

func TestGetProductInvalid(t *testing.T) {
	ts := setupTestServer()
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/999", nil)
	ts.GetProduct(r, req, 999)

	if r.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", r.Code)
	}
}

func TestPlaceOrderWithValidPromo(t *testing.T) {
	ts := setupTestServer()

	order := api.OrderReq{
		CouponCode: ptrString("HAPPYHRS"),
		Items: []struct {
			ProductId string `json:"productId"`
			Quantity  int    `json:"quantity"`
		}{
			{ProductId: "1", Quantity: 2},
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(order)
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/order", buf)
	ts.PlaceOrder(r, req)

	if r.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", r.Code)
	}

	var res api.Order
	_ = json.NewDecoder(r.Body).Decode(&res)
	if res.Discounts == nil || *res.Discounts == 0 {
		t.Errorf("Expected discount to be applied, but got discount of %v", res.Discounts)
	}
}

func TestPlaceOrderWithInvalidPromo(t *testing.T) {
	ts := setupTestServer()

	order := api.OrderReq{
		CouponCode: ptrString("INVALID123"),
		Items: []struct {
			ProductId string `json:"productId"`
			Quantity  int    `json:"quantity"`
		}{
			{ProductId: "1", Quantity: 2},
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(order)
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/order", buf)
	ts.PlaceOrder(r, req)

	if r.Code != http.StatusBadRequest {
		t.Errorf("expected %d Bad Request, got %d", http.StatusBadRequest, r.Code)
	}

	expectedErr := "Invalid promo code"
	if body := r.Body.String(); !strings.Contains(body, expectedErr) {
		t.Errorf("expected response body to contain '%s', got '%s'", expectedErr, body)
	}
}

func TestPlaceOrderInvalidProduct(t *testing.T) {
	ts := setupTestServer()

	order := api.OrderReq{
		Items: []struct {
			ProductId string `json:"productId"`
			Quantity  int    `json:"quantity"`
		}{
			{ProductId: "999", Quantity: 2},
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(order)
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/order", buf)
	ts.PlaceOrder(r, req)

	if r.Code != http.StatusBadRequest {
		t.Errorf("expected %d Bad Request, got %d", http.StatusBadRequest, r.Code)
	}

	expectedErr := "Invalid product ID: 999"
	if body := r.Body.String(); !strings.Contains(body, expectedErr) {
		t.Errorf("expected response body to contain '%s', got '%s'", expectedErr, body)
	}
}
