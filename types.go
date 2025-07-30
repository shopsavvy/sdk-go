package shopsavvy

// APIResponse represents a response from the ShopSavvy API
type APIResponse[T any] struct {
	Success          bool   `json:"success"`
	Data             T      `json:"data"`
	Message          string `json:"message,omitempty"`
	CreditsUsed      *int   `json:"credits_used,omitempty"`
	CreditsRemaining *int   `json:"credits_remaining,omitempty"`
}

// ProductDetails represents detailed product information
type ProductDetails struct {
	ProductID   string            `json:"product_id"`
	Name        string            `json:"name"`
	Brand       *string           `json:"brand,omitempty"`
	Category    *string           `json:"category,omitempty"`
	ImageURL    *string           `json:"image_url,omitempty"`
	Barcode     *string           `json:"barcode,omitempty"`
	ASIN        *string           `json:"asin,omitempty"`
	Model       *string           `json:"model,omitempty"`
	MPN         *string           `json:"mpn,omitempty"`
	Description *string           `json:"description,omitempty"`
	Identifiers map[string]string `json:"identifiers,omitempty"`
}

// Offer represents a product offer from a retailer
type Offer struct {
	OfferID      string  `json:"offer_id"`
	Retailer     string  `json:"retailer"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Availability string  `json:"availability"`
	Condition    string  `json:"condition"`
	URL          string  `json:"url"`
	Shipping     *float64 `json:"shipping,omitempty"`
	LastUpdated  string  `json:"last_updated"`
}

// PriceHistoryEntry represents a single price point in history
type PriceHistoryEntry struct {
	Date         string `json:"date"`
	Price        float64 `json:"price"`
	Availability string `json:"availability"`
}

// OfferWithHistory represents an offer with price history
type OfferWithHistory struct {
	Offer
	PriceHistory []PriceHistoryEntry `json:"price_history"`
}

// ScheduledProduct represents a product that is scheduled for monitoring
type ScheduledProduct struct {
	ProductID     string  `json:"product_id"`
	Identifier    string  `json:"identifier"`
	Frequency     string  `json:"frequency"`
	Retailer      *string `json:"retailer,omitempty"`
	CreatedAt     string  `json:"created_at"`
	LastRefreshed *string `json:"last_refreshed,omitempty"`
}

// UsageInfo represents API usage and credit information
type UsageInfo struct {
	CreditsUsed          int    `json:"credits_used"`
	CreditsRemaining     int    `json:"credits_remaining"`
	CreditsTotal         int    `json:"credits_total"`
	BillingPeriodStart   string `json:"billing_period_start"`
	BillingPeriodEnd     string `json:"billing_period_end"`
	PlanName             string `json:"plan_name"`
}

// ScheduleResponse represents the response from scheduling a product
type ScheduleResponse struct {
	Scheduled bool   `json:"scheduled"`
	ProductID string `json:"product_id"`
}

// ScheduleBatchResponse represents the response from batch scheduling
type ScheduleBatchResponse struct {
	Identifier string `json:"identifier"`
	Scheduled  bool   `json:"scheduled"`
	ProductID  string `json:"product_id"`
}

// RemoveResponse represents the response from removing a product from schedule
type RemoveResponse struct {
	Removed bool `json:"removed"`
}

// RemoveBatchResponse represents the response from batch removal
type RemoveBatchResponse struct {
	Identifier string `json:"identifier"`
	Removed    bool   `json:"removed"`
}

// APIErrorResponse represents an error response from the API
type APIErrorResponse struct {
	Error string `json:"error"`
}