package shopsavvy

// APIMeta contains credit usage info from the API response
type APIMeta struct {
	CreditsUsed        int  `json:"credits_used"`
	CreditsRemaining   int  `json:"credits_remaining"`
	RateLimitRemaining *int `json:"rate_limit_remaining,omitempty"`
}

// APIResponse represents a response from the ShopSavvy API
type APIResponse[T any] struct {
	Success bool    `json:"success"`
	Data    T       `json:"data"`
	Message string  `json:"message,omitempty"`
	Meta    *APIMeta `json:"meta,omitempty"`
}

// CreditsUsed returns the credits used from the meta object
func (r *APIResponse[T]) CreditsUsed() int {
	if r.Meta != nil {
		return r.Meta.CreditsUsed
	}
	return 0
}

// CreditsRemaining returns the credits remaining from the meta object
func (r *APIResponse[T]) CreditsRemaining() int {
	if r.Meta != nil {
		return r.Meta.CreditsRemaining
	}
	return 0
}

// ProductDetails represents detailed product information
type ProductDetails struct {
	Title     string   `json:"title"`
	ShopSavvy string   `json:"shopsavvy"`
	Brand     *string  `json:"brand,omitempty"`
	Category  *string  `json:"category,omitempty"`
	Images    []string `json:"images,omitempty"`
	Barcode   *string  `json:"barcode,omitempty"`
	Amazon    *string  `json:"amazon,omitempty"`
	Model     *string  `json:"model,omitempty"`
	MPN       *string  `json:"mpn,omitempty"`
	Color       *string            `json:"color,omitempty"`
	TitleShort  *string            `json:"title_short,omitempty"`
	Slug        *string            `json:"slug,omitempty"`
	Description *string            `json:"description,omitempty"`
	Categories  []string           `json:"categories,omitempty"`
	Attributes  map[string]string  `json:"attributes,omitempty"`
	Rating      *ProductRating     `json:"rating,omitempty"`
	Score       *ProductScore      `json:"score,omitempty"`
	Keywords    []string           `json:"keywords,omitempty"`
	Identifiers map[string]any     `json:"identifiers,omitempty"`
}

// Name returns the product title (deprecated alias)
func (p *ProductDetails) Name() string {
	return p.Title
}

// ProductID returns the ShopSavvy ID (deprecated alias)
func (p *ProductDetails) ProductID() string {
	return p.ShopSavvy
}

// ASIN returns the Amazon ASIN (deprecated alias)
func (p *ProductDetails) ASIN() *string {
	return p.Amazon
}

// ImageURL returns the first image URL (deprecated alias)
func (p *ProductDetails) ImageURL() *string {
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

// Offer represents a product offer from a retailer
type Offer struct {
	ID           string              `json:"id"`
	Retailer     *string             `json:"retailer,omitempty"`
	Price        *float64            `json:"price,omitempty"`
	Currency     *string             `json:"currency,omitempty"`
	Availability *string             `json:"availability,omitempty"`
	Condition    *string             `json:"condition,omitempty"`
	URL          *string             `json:"URL,omitempty"`
	Seller       *string             `json:"seller,omitempty"`
	Timestamp    *string             `json:"timestamp,omitempty"`
	History      []PriceHistoryEntry `json:"history,omitempty"`
}

// OfferID returns the offer ID (deprecated alias)
func (o *Offer) OfferID() string {
	return o.ID
}

// OfferURL returns the offer URL (deprecated alias)
func (o *Offer) OfferURL() *string {
	return o.URL
}

// LastUpdated returns the timestamp (deprecated alias)
func (o *Offer) LastUpdated() *string {
	return o.Timestamp
}

// ProductWithOffers represents a product with its current offers
type ProductWithOffers struct {
	ProductDetails
	Offers []Offer `json:"offers"`
}

// PriceHistoryEntry represents a single price point in history
type PriceHistoryEntry struct {
	Date         string  `json:"date"`
	Price        float64 `json:"price"`
	Availability string  `json:"availability"`
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

// UsagePeriod represents the current billing period details
type UsagePeriod struct {
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	CreditsUsed      int    `json:"credits_used"`
	CreditsLimit     int    `json:"credits_limit"`
	CreditsRemaining int    `json:"credits_remaining"`
	RequestsMade     int    `json:"requests_made"`
}

// UsageInfo represents API usage and credit information
type UsageInfo struct {
	CurrentPeriod   UsagePeriod `json:"current_period"`
	UsagePercentage float64     `json:"usage_percentage"`
}

// Deprecated accessors for backward compatibility

// GetCreditsUsed returns credits used (deprecated, use CurrentPeriod.CreditsUsed)
func (u *UsageInfo) GetCreditsUsed() int {
	return u.CurrentPeriod.CreditsUsed
}

// GetCreditsRemaining returns credits remaining (deprecated, use CurrentPeriod.CreditsRemaining)
func (u *UsageInfo) GetCreditsRemaining() int {
	return u.CurrentPeriod.CreditsRemaining
}

// GetCreditsTotal returns credits limit (deprecated, use CurrentPeriod.CreditsLimit)
func (u *UsageInfo) GetCreditsTotal() int {
	return u.CurrentPeriod.CreditsLimit
}

// GetBillingPeriodStart returns billing period start (deprecated, use CurrentPeriod.StartDate)
func (u *UsageInfo) GetBillingPeriodStart() string {
	return u.CurrentPeriod.StartDate
}

// GetBillingPeriodEnd returns billing period end (deprecated, use CurrentPeriod.EndDate)
func (u *UsageInfo) GetBillingPeriodEnd() string {
	return u.CurrentPeriod.EndDate
}

// PaginationInfo represents pagination metadata for search results
type PaginationInfo struct {
	Total    int `json:"total"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	Returned int `json:"returned"`
}

// ProductRating represents an aggregated customer rating
type ProductRating struct {
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

// ProductScore represents unified expert quality scores on a 0-1 scale
// (multiply by 10 or 100 for display). Aspects is keyed by free-form aspect
// names drawn from the product's professional reviews (e.g. "value",
// "features", "reliability", "quality") — the set of aspects varies per product.
type ProductScore struct {
	Overall      *float64           `json:"overall,omitempty"`
	Customer     *float64           `json:"customer,omitempty"`
	Professional *float64           `json:"professional,omitempty"`
	Aspects      map[string]float64 `json:"aspects,omitempty"`
}

// Deal represents a shopping deal with expert grading
type Deal struct {
	Path         string          `json:"path"`
	Title        string          `json:"title"`
	Subtitle     *string         `json:"subtitle,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Emoji        *string         `json:"emoji,omitempty"`
	Grade        DealGrade       `json:"grade"`
	Pricing      DealPricing     `json:"pricing"`
	Retailer     DealRetailer    `json:"retailer"`
	Product      *string         `json:"product,omitempty"`
	URL          string          `json:"url"`
	Image        map[string]any  `json:"image,omitempty"`
	Votes        DealVotes       `json:"votes"`
	CommentCount int             `json:"comment_count"`
	Tags         []DealTag       `json:"tags,omitempty"`
	ExpiresAt    *string         `json:"expires_at,omitempty"`
	CreatedAt    string          `json:"created_at"`
}

type DealGrade struct {
	Letter        string   `json:"letter"`
	Suffix        *string  `json:"suffix,omitempty"`
	Value         float64  `json:"value"`
	Justification *string  `json:"justification,omitempty"`
}

type DealPricing struct {
	Current  float64  `json:"current"`
	Original *float64 `json:"original,omitempty"`
	Currency string   `json:"currency"`
}

type DealRetailer struct {
	Name string `json:"name"`
}

type DealVotes struct {
	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
	Score     int `json:"score"`
}

type DealTag struct {
	Slug    string `json:"slug"`
	Display string `json:"display"`
}

// DealsResponse represents the response from the deals endpoint
type DealsResponse struct {
	Success    bool          `json:"success"`
	Deals      []Deal        `json:"deals"`
	Pagination DealsPagination `json:"pagination"`
	Meta       *APIMeta      `json:"meta,omitempty"`
}

type DealsPagination struct {
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
}

// DealsOptions configures the deals request
type DealsOptions struct {
	Sort     string
	Limit    int
	Offset   int
	Category string
	Retailer string
	Tag      string
	MinPrice float64
	MaxPrice float64
	Grade    string
}

// TLDRReview represents an expert product review summary
type TLDRReview struct {
	Slug       string        `json:"slug"`
	Headline   string        `json:"headline"`
	Pros       []string      `json:"pros"`
	Cons       []string      `json:"cons"`
	BottomLine string        `json:"bottom_line"`
	Scores     *ProductScore `json:"scores,omitempty"`
}

// ReviewResponse represents the response from the reviews endpoint
type ReviewResponse struct {
	Success bool               `json:"success"`
	Product map[string]string  `json:"product"`
	Review  *TLDRReview        `json:"review"`
	Meta    *APIMeta           `json:"meta,omitempty"`
}

// ProductSearchResult represents the response from a product search
type ProductSearchResult struct {
	Success    bool             `json:"success"`
	Data       []ProductDetails `json:"data"`
	Pagination *PaginationInfo  `json:"pagination,omitempty"`
	Meta       *APIMeta         `json:"meta,omitempty"`
}

// CreditsUsed returns the credits used from the meta object
func (r *ProductSearchResult) CreditsUsed() int {
	if r.Meta != nil {
		return r.Meta.CreditsUsed
	}
	return 0
}

// CreditsRemaining returns the credits remaining from the meta object
func (r *ProductSearchResult) CreditsRemaining() int {
	if r.Meta != nil {
		return r.Meta.CreditsRemaining
	}
	return 0
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
