// Package shopsavvy provides a Go client for the ShopSavvy Data API
//
// This SDK provides a convenient interface to interact with the ShopSavvy Data API,
// allowing you to access product data, pricing information, and price history
// across thousands of retailers and millions of products.
//
// Example usage:
//
//	client, err := shopsavvy.NewClient("ss_live_your_api_key_here")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	product, err := client.GetProductDetails("012345678901")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(product.Data[0].Title)
package shopsavvy

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Version is the current SDK version
const Version = "1.0.1"

// Client represents the ShopSavvy Data API client
type Client struct {
	config *Config
	client *resty.Client
}

// Config holds the configuration for the ShopSavvy API client
type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// NewClient creates a new ShopSavvy Data API client with the given API key
func NewClient(apiKey string, options ...Option) (*Client, error) {
	config := &Config{
		APIKey:  apiKey,
		BaseURL: "https://api.shopsavvy.com/v1",
		Timeout: 30 * time.Second,
	}

	// Apply options
	for _, option := range options {
		option(config)
	}

	// Validate API key
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required. Get one at https://shopsavvy.com/data")
	}

	// Validate API key format
	matched, _ := regexp.MatchString(`^ss_(live|test)_[a-zA-Z0-9]+$`, config.APIKey)
	if !matched {
		return nil, fmt.Errorf("invalid API key format. API keys should start with ss_live_ or ss_test_")
	}

	// Create HTTP client
	client := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetHeader("Authorization", "Bearer "+config.APIKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "ShopSavvy-Go-SDK/"+Version).
		SetError(&APIErrorResponse{}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			if resp.IsError() {
				return handleErrorResponse(resp)
			}
			return nil
		})

	return &Client{
		config: config,
		client: client,
	}, nil
}

// Option is a functional option for configuring the client
type Option func(*Config)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets a custom timeout for API requests
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// Close closes the HTTP client and releases resources
func (c *Client) Close() {
	// resty client doesn't need explicit closing, but we provide this for consistency
}

// SearchProducts searches for products by keyword
func (c *Client) SearchProducts(query string, limit, offset int) (*ProductSearchResult, error) {
	params := map[string]string{
		"q": query,
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	if offset > 0 {
		params["offset"] = fmt.Sprintf("%d", offset)
	}

	var response ProductSearchResult
	_, err := c.client.R().
		SetQueryParams(params).
		SetResult(&response).
		Get("/products/search")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProductDetails looks up product details by identifier
func (c *Client) GetProductDetails(identifier string, format ...string) (*APIResponse[[]ProductDetails], error) {
	params := map[string]string{
		"ids": identifier,
	}
	if len(format) > 0 && format[0] != "" {
		params["format"] = format[0]
	}

	var response APIResponse[[]ProductDetails]
	_, err := c.client.R().
		SetQueryParams(params).
		SetResult(&response).
		Get("/products")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProductDetailsBatch looks up details for multiple products
func (c *Client) GetProductDetailsBatch(identifiers []string, format ...string) (*APIResponse[[]ProductDetails], error) {
	params := map[string]string{
		"ids": strings.Join(identifiers, ","),
	}
	if len(format) > 0 && format[0] != "" {
		params["format"] = format[0]
	}

	var response APIResponse[[]ProductDetails]
	_, err := c.client.R().
		SetQueryParams(params).
		SetResult(&response).
		Get("/products")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCurrentOffers gets current offers for a product
func (c *Client) GetCurrentOffers(identifier string, retailer string, format ...string) (*APIResponse[[]ProductWithOffers], error) {
	params := map[string]string{
		"ids": identifier,
	}
	if retailer != "" {
		params["retailer"] = retailer
	}
	if len(format) > 0 && format[0] != "" {
		params["format"] = format[0]
	}

	var response APIResponse[[]ProductWithOffers]
	_, err := c.client.R().
		SetQueryParams(params).
		SetResult(&response).
		Get("/products/offers")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCurrentOffersBatch gets current offers for multiple products
func (c *Client) GetCurrentOffersBatch(identifiers []string, retailer string, format ...string) (*APIResponse[[]ProductWithOffers], error) {
	params := map[string]string{
		"ids": strings.Join(identifiers, ","),
	}
	if retailer != "" {
		params["retailer"] = retailer
	}
	if len(format) > 0 && format[0] != "" {
		params["format"] = format[0]
	}

	var response APIResponse[[]ProductWithOffers]
	_, err := c.client.R().
		SetQueryParams(params).
		SetResult(&response).
		Get("/products/offers")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPriceHistory gets price history for a product
func (c *Client) GetPriceHistory(identifier, startDate, endDate string, retailer string, format ...string) (*APIResponse[[]OfferWithHistory], error) {
	params := map[string]string{
		"ids":        identifier,
		"start_date": startDate,
		"end_date":   endDate,
	}
	if retailer != "" {
		params["retailer"] = retailer
	}
	if len(format) > 0 && format[0] != "" {
		params["format"] = format[0]
	}

	var response APIResponse[[]OfferWithHistory]
	_, err := c.client.R().
		SetQueryParams(params).
		SetResult(&response).
		Get("/products/offers/history")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScheduleProductMonitoring schedules product monitoring
func (c *Client) ScheduleProductMonitoring(identifier, frequency string, retailer ...string) (*APIResponse[ScheduleResponse], error) {
	body := map[string]interface{}{
		"identifier": identifier,
		"frequency":  frequency,
	}
	if len(retailer) > 0 && retailer[0] != "" {
		body["retailer"] = retailer[0]
	}

	var response APIResponse[ScheduleResponse]
	_, err := c.client.R().
		SetBody(body).
		SetResult(&response).
		Post("/products/schedule")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScheduleProductMonitoringBatch schedules monitoring for multiple products
func (c *Client) ScheduleProductMonitoringBatch(identifiers []string, frequency string, retailer ...string) (*APIResponse[[]ScheduleBatchResponse], error) {
	body := map[string]interface{}{
		"identifiers": strings.Join(identifiers, ","),
		"frequency":   frequency,
	}
	if len(retailer) > 0 && retailer[0] != "" {
		body["retailer"] = retailer[0]
	}

	var response APIResponse[[]ScheduleBatchResponse]
	_, err := c.client.R().
		SetBody(body).
		SetResult(&response).
		Post("/products/schedule")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetScheduledProducts gets all scheduled products
func (c *Client) GetScheduledProducts() (*APIResponse[[]ScheduledProduct], error) {
	var response APIResponse[[]ScheduledProduct]
	_, err := c.client.R().
		SetResult(&response).
		Get("/products/scheduled")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveProductFromSchedule removes a product from monitoring schedule
func (c *Client) RemoveProductFromSchedule(identifier string) (*APIResponse[RemoveResponse], error) {
	body := map[string]interface{}{
		"identifier": identifier,
	}

	var response APIResponse[RemoveResponse]
	_, err := c.client.R().
		SetBody(body).
		SetResult(&response).
		Delete("/products/schedule")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveProductsFromSchedule removes multiple products from monitoring schedule
func (c *Client) RemoveProductsFromSchedule(identifiers []string) (*APIResponse[[]RemoveBatchResponse], error) {
	body := map[string]interface{}{
		"identifiers": strings.Join(identifiers, ","),
	}

	var response APIResponse[[]RemoveBatchResponse]
	_, err := c.client.R().
		SetBody(body).
		SetResult(&response).
		Delete("/products/schedule")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUsage gets API usage information
func (c *Client) GetUsage() (*APIResponse[UsageInfo], error) {
	var response APIResponse[UsageInfo]
	_, err := c.client.R().
		SetResult(&response).
		Get("/usage")

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// handleErrorResponse converts HTTP error responses to appropriate Go errors
func handleErrorResponse(resp *resty.Response) error {
	statusCode := resp.StatusCode()

	var errorMsg string
	if errorResp, ok := resp.Error().(*APIErrorResponse); ok && errorResp.Error != "" {
		errorMsg = errorResp.Error
	} else {
		errorMsg = fmt.Sprintf("HTTP %d: %s", statusCode, http.StatusText(statusCode))
	}

	switch statusCode {
	case 401:
		return &AuthenticationError{
			Message:    "Authentication failed. Check your API key.",
			StatusCode: statusCode,
		}
	case 404:
		return &NotFoundError{
			Message:    "Resource not found",
			StatusCode: statusCode,
		}
	case 422:
		return &ValidationError{
			Message:    "Request validation failed. Check your parameters.",
			StatusCode: statusCode,
		}
	case 429:
		return &RateLimitError{
			Message:    "Rate limit exceeded. Please slow down your requests.",
			StatusCode: statusCode,
		}
	default:
		return &APIError{
			Message:    errorMsg,
			StatusCode: statusCode,
		}
	}
}
