// impl/server.go
package impl

import (
	"backend-challenge/api"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	// DiscountPercent defines the % discount for valid promo codes.
	// This is a configurable business rule — not defined in the OpenAPI spec.
	DiscountPercent = 10.0
)

// Server implements api.ServerInterface
// It defines handlers for each endpoint described in the OpenAPI spec.
type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// ListProducts handles GET /product
func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {
	products := ListAllProducts()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetProduct handles GET /product/{productId}
func (s *Server) GetProduct(w http.ResponseWriter, r *http.Request, productId int64) {
	p := GetProductByID(strconv.FormatInt(productId, 10))
	if p == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// PlaceOrder handles POST /order
func (s *Server) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req api.OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var orderItems []api.Product
	total := float32(0.0)
	for _, item := range req.Items {
		p := GetProductByID(item.ProductId)
		if p == nil {
			http.Error(w, "Invalid product ID: "+item.ProductId, http.StatusBadRequest)
			return
		}
		if item.Quantity <= 0 {
			http.Error(w, "Quantity must be > 0", http.StatusBadRequest)
			return
		}
		total += *p.Price * float32(item.Quantity)
		orderItems = append(orderItems, *p)
	}

	// Promo code validation
	discount := float32(0.0)
	if req.CouponCode != nil {
		code := strings.TrimSpace(*req.CouponCode)
		if !IsPromoCodeValid(code) {
			http.Error(w, "Invalid promo code", http.StatusBadRequest)
			return
		}

		// Apply discount. Assuming static % logic here, production would have a very well defined logic for this
		discount = float32(DiscountPercent/100.0) * total
		total -= discount
	}

	response := api.Order{
		Id: ptrString(uuid.New().String()),
		Items: &[]struct {
			ProductId *string `json:"productId,omitempty"`
			Quantity  *int    `json:"quantity,omitempty"`
		}{},
		Products:  &orderItems,
		Total:     &total,
		Discounts: &discount,
	}
	for _, item := range req.Items {
		pid := item.ProductId
		qty := item.Quantity
		*response.Items = append(*response.Items, struct {
			ProductId *string `json:"productId,omitempty"`
			Quantity  *int    `json:"quantity,omitempty"`
		}{ProductId: &pid, Quantity: &qty})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
