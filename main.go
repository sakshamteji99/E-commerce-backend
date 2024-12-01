package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

type Admin struct {
	IsSet  bool
	Secret string
}
type Product struct {
	ID        int
	Name      string
	Stock     int
	SalePrice float64
}
type Coupon struct {
	ID               int
	Category         string
	MinCartValue     float64
	MaxCartValue     float64
	DiscountPercent  float64
	MaxDiscountValue float64
	Description      string
}
type Customer struct {
	ID         int
	Wallet     float64
	Points     int
	BirthMonth time.Month
	OrderCount int
}

var (
	admin             = Admin{IsSet: false, Secret: "Admin123"}
	products          = make(map[int]Product)
	coupons           = make(map[int]Coupon)
	customers         = make(map[int]Customer)
	couponCategories  = []string{}
	productIDCounter  = 1
	couponIDCounter   = 1
	customerIDCounter = 1
)

const (
	freeDeliveryThreshold = 1000.0
	extraDiscountPercent  = 5.0
)

// calculate discounts
func calculateDiscount(cartValue, discountPercent, maxDiscount float64) float64 {
	discount := (cartValue * discountPercent) / 100
	if discount > maxDiscount {
		return maxDiscount
	}
	return discount
}

// check if today is a special sale day
func isSpecialSaleDay() bool {
	today := time.Now()
	month, day := today.Month(), today.Day()
	return (month == time.December && day == 25) || (month == time.August && day == 15) // Christmas sale and Independence day sale

}

// Set Admin API
func SetAdmin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Secret string `json:"secret"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if admin.IsSet {
		http.Error(w, "Admin already set", http.StatusForbidden)
		return
	}
	if req.Secret == admin.Secret {
		admin.IsSet = true
		w.Write([]byte(`{"message": "Admin set successfully"}`))
	} else {
		http.Error(w, "Invalid secret", http.StatusForbidden)
	}
}

// Add or Remove Product API
func AddOrRemoveProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Secret    string  `json:"secret"`
		Action    string  `json:"action"`
		Name      string  `json:"name"`
		Quantity  int     `json:"quantity"`
		SalePrice float64 `json:"sale_price"`
		ProductID int     `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Secret != admin.Secret {
		http.Error(w, "Invalid admin secret", http.StatusForbidden)
		return
	}
	if req.Action == "add" {
		product := Product{
			ID:        productIDCounter,
			Name:      req.Name,
			Stock:     req.Quantity,
			SalePrice: req.SalePrice,
		}
		products[productIDCounter] = product
		productIDCounter++
		w.Write([]byte(fmt.Sprintf(`{"message": "Product added", "product_id": %d}`, product.ID)))
	} else if req.Action == "remove" {
		if _, exists := products[req.ProductID]; exists {
			delete(products, req.ProductID)
			w.Write([]byte(`{"message": "Product removed"}`))
		} else {
			http.Error(w, "Product not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

// Add or Remove Coupon Category API
func AddOrRemoveCouponCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Action == "add" {
		couponCategories = append(couponCategories, req.Name)
		w.Write([]byte(`{"message": "Coupon category added"}`))
	} else if req.Action == "remove" {
		for i, category := range couponCategories {
			if strings.EqualFold(category, req.Name) {
				couponCategories = append(couponCategories[:i], couponCategories[i+1:]...)
				w.Write([]byte(`{"message": "Coupon category removed"}`))
				return
			}
		}
		http.Error(w, "Category not found", http.StatusNotFound)
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

// Add or Remove Coupon API
func AddOrRemoveCoupon(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Secret      string  `json:"secret"`
		Action      string  `json:"action"` // "add" or "remove"
		Category    string  `json:"category"`
		MinValue    float64 `json:"min_value"`
		MaxValue    float64 `json:"max_value"`
		Discount    float64 `json:"discount_percent"`
		MaxDiscount float64 `json:"max_discount"`
		CouponID    int     `json:"coupon_id"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Secret != admin.Secret {
		http.Error(w, "Invalid admin secret", http.StatusForbidden)
		return
	}
	if req.Action == "add" {
		coupon := Coupon{
			ID:               couponIDCounter,
			Category:         req.Category,
			MinCartValue:     req.MinValue,
			MaxCartValue:     req.MaxValue,
			DiscountPercent:  req.Discount,
			MaxDiscountValue: req.MaxDiscount,
			Description:      req.Description,
		}
		coupons[couponIDCounter] = coupon
		couponIDCounter++
		w.Write([]byte(fmt.Sprintf(`{"message": "Coupon added", "coupon_id": %d}`, coupon.ID)))
	} else if req.Action == "remove" {
		if _, exists := coupons[req.CouponID]; exists {
			delete(coupons, req.CouponID)
			w.Write([]byte(`{"message": "Coupon removed"}`))
		} else {
			http.Error(w, "Coupon not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

// Check Discounts API
func CheckDiscounts(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CartValue  float64 `json:"cart_value"`
		CustomerID int     `json:"customer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	customer, exists := customers[req.CustomerID]
	if !exists {
		http.Error(w, "Customer not found", http.StatusBadRequest)
		return
	}
	var discountDetails []map[string]interface{}
	//applicable coupons
	for _, coupon := range coupons {
		if req.CartValue >= coupon.MinCartValue && req.CartValue <= coupon.MaxCartValue {
			discount := calculateDiscount(req.CartValue, coupon.DiscountPercent, coupon.MaxDiscountValue)
			discountDetails = append(discountDetails, map[string]interface{}{
				"type":           "Coupon Discount",
				"coupon_id":      coupon.ID,
				"description":    coupon.Description,
				"discount_value": discount,
			})
		}
	}
	// Special sale day discount
	if isSpecialSaleDay() {
		specialSaleDiscount := calculateDiscount(req.CartValue, extraDiscountPercent, math.MaxFloat64)
		discountDetails = append(discountDetails, map[string]interface{}{
			"type":           "Special Sale Discount",
			"description":    "Extra discount for special sale days",
			"discount_value": specialSaleDiscount,
		})
	}
	// Birthday discount
	if time.Now().Month() == customer.BirthMonth {
		birthdayDiscount := calculateDiscount(req.CartValue, extraDiscountPercent, math.MaxFloat64)
		discountDetails = append(discountDetails, map[string]interface{}{
			"type":           "Birthday Discount",
			"description":    "Extra discount for your birthday month",
			"discount_value": birthdayDiscount,
		})
	}
	// Anniversary discount
	today := time.Now()
	if today.Month() == time.December && today.Day() == 1 {
		anniversaryDiscount := calculateDiscount(req.CartValue, 5.0, math.MaxFloat64)
		discountDetails = append(discountDetails, map[string]interface{}{
			"type":           "Anniversary Discount",
			"description":    "Flat 5% discount on company anniversary",
			"discount_value": anniversaryDiscount,
		})
	}
	// Calculate maximum cash discount and points equivalent
	totalCashDiscount := 0.0
	for _, discount := range discountDetails {
		totalCashDiscount += discount["discount_value"].(float64)
	}
	pointsEquivalent := int(totalCashDiscount * 1.2)
	// Add summary
	discountDetails = append(discountDetails, map[string]interface{}{
		"type":                "Summary",
		"total_cash_discount": totalCashDiscount,
		"points_equivalent":   pointsEquivalent,
		"description":         "You can choose between total cash discount or points equivalent.",
	})
	respJSON, _ := json.Marshal(discountDetails)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

// Apply Discount API with Points Option
func ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CartValue  float64 `json:"cart_value"`
		CouponID   int     `json:"coupon_id"`
		CustomerID int     `json:"customer_id"`
		Option     string  `json:"option"` // "cash" or "points"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	coupon, exists := coupons[req.CouponID]
	if !exists {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}
	customer, exists := customers[req.CustomerID]
	if !exists {
		http.Error(w, "Customer not found", http.StatusBadRequest)
		return
	}
	// Calculate cash discount
	cashDiscount := calculateDiscount(req.CartValue, coupon.DiscountPercent, coupon.MaxDiscountValue)
	// Calculate special sale and birthday discounts
	if isSpecialSaleDay() {
		cashDiscount += calculateDiscount(req.CartValue, extraDiscountPercent, math.MaxFloat64)
	}
	if time.Now().Month() == customer.BirthMonth {
		cashDiscount += calculateDiscount(req.CartValue, extraDiscountPercent, math.MaxFloat64)
	}
	// Option 1: Cash Discount
	finalCartValue := req.CartValue - cashDiscount
	if finalCartValue < 0 {
		finalCartValue = 0
	}
	// Option 2: Points (1.2x cash discount)
	pointsEquivalent := int(cashDiscount * 1.2)
	if req.Option == "cash" {
		// Apply cash discount
		resp := map[string]interface{}{
			"original_cart_value": req.CartValue,
			"final_cart_value":    finalCartValue,
			"total_discount":      cashDiscount,
			"points_earned":       0,
			"message":             "Cash discount applied successfully",
		}
		respJSON, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respJSON)
	} else if req.Option == "points" {
		// Apply points
		customer.Points += pointsEquivalent
		customers[req.CustomerID] = customer

		resp := map[string]interface{}{
			"original_cart_value": req.CartValue,
			"final_cart_value":    req.CartValue,
			"total_discount":      0,
			"points_earned":       pointsEquivalent,
			"new_total_points":    customer.Points,
			"message":             "Points credited successfully",
		}
		respJSON, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respJSON)
	} else {
		http.Error(w, "Invalid option. Choose 'cash' or 'points'.", http.StatusBadRequest)
	}
}

// Place Order API
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID int     `json:"customer_id"`
		ProductID  int     `json:"product_id"`
		Quantity   int     `json:"quantity"`
		CartValue  float64 `json:"cart_value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	product, exists := products[req.ProductID]
	if !exists || product.Stock < req.Quantity {
		http.Error(w, "Product not available or insufficient stock", http.StatusBadRequest)
		return
	}
	product.Stock -= req.Quantity
	products[req.ProductID] = product
	customer, exists := customers[req.CustomerID]
	if !exists {
		http.Error(w, "Customer not found", http.StatusBadRequest)
		return
	}
	// Determine if free delivery is applicable
	freeDelivery := false
	if customer.OrderCount < 5 && req.CartValue >= freeDeliveryThreshold {
		freeDelivery = true
	}
	// Update customer order count and points
	customer.OrderCount++
	pointsEarned := int(product.SalePrice * float64(req.Quantity) / 10)
	customer.Points += pointsEarned
	customers[req.CustomerID] = customer
	resp := map[string]interface{}{
		"message":             "Order placed successfully",
		"free_delivery":       freeDelivery,
		"points_earned":       pointsEarned,
		"remaining_stock":     product.Stock,
		"updated_order_count": customer.OrderCount,
	}
	respJSON, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

// Add Money to Wallet API
func AddMoneyToWallet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID int     `json:"customer_id"`
		Amount     float64 `json:"amount"`
		UsePoints  bool    `json:"use_points"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	customer, exists := customers[req.CustomerID]
	if !exists {
		http.Error(w, "Customer not found", http.StatusBadRequest)
		return
	}
	if req.UsePoints {
		amountFromPoints := float64(customer.Points)
		customer.Wallet += amountFromPoints
		customer.Points = 0
	} else {
		customer.Wallet += req.Amount
	}
	customers[req.CustomerID] = customer
	resp := map[string]interface{}{
		"message":          "Wallet updated successfully",
		"new_wallet":       customer.Wallet,
		"remaining_points": customer.Points,
	}
	respJSON, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func main() {
	http.HandleFunc("/admin/set-admin", SetAdmin)
	http.HandleFunc("/admin/add-or-remove/product", AddOrRemoveProduct)
	http.HandleFunc("/admin/add-or-remove/coupon/category", AddOrRemoveCouponCategory)
	http.HandleFunc("/admin/add-or-remove/coupon/", AddOrRemoveCoupon)
	http.HandleFunc("/customer/check-discounts", CheckDiscounts)
	http.HandleFunc("/customer/apply-discount", ApplyDiscount)
	http.HandleFunc("/customer/place-order", PlaceOrder)
	http.HandleFunc("/customer/add-money-to-wallet", AddMoneyToWallet)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
