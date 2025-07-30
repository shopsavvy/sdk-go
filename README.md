# ShopSavvy Go SDK

Official Go SDK for the ShopSavvy Data API - Access product data, pricing information, and price history across thousands of retailers and millions of products.

## Installation

```bash
go get github.com/shopsavvy/sdk-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/shopsavvy/sdk-go"
)

func main() {
    // Create a new client
    client, err := shopsavvy.NewClient("ss_live_your_api_key_here")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Look up a product by barcode
    product, err := client.GetProductDetails("012345678901")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Product: %s\n", product.Data.Name)
    if product.Data.Brand != nil {
        fmt.Printf("Brand: %s\n", *product.Data.Brand)
    }
    
    // Get current offers
    offers, err := client.GetCurrentOffers("012345678901", "")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d offers:\n", len(offers.Data))
    for _, offer := range offers.Data {
        fmt.Printf("  %s: $%.2f (%s)\n", offer.Retailer, offer.Price, offer.Availability)
    }
}
```

## Configuration Options

You can customize the client behavior using functional options:

```go
import "time"

client, err := shopsavvy.NewClient(
    "ss_live_your_api_key_here",
    shopsavvy.WithTimeout(60*time.Second),
    shopsavvy.WithBaseURL("https://api.shopsavvy.com/v1"),
)
```

## API Methods

### Product Details

```go
// Single product
product, err := client.GetProductDetails("012345678901")

// Multiple products
products, err := client.GetProductDetailsBatch([]string{"012345678901", "B08N5WRWNW"})
```

### Current Offers

```go
// All retailers
offers, err := client.GetCurrentOffers("012345678901", "")

// Specific retailer
offers, err := client.GetCurrentOffers("012345678901", "amazon")

// Multiple products
offersBatch, err := client.GetCurrentOffersBatch([]string{"012345678901", "B08N5WRWNW"}, "")
```

### Price History

```go
history, err := client.GetPriceHistory("012345678901", "2024-01-01", "2024-01-31", "")

// With specific retailer
history, err := client.GetPriceHistory("012345678901", "2024-01-01", "2024-01-31", "amazon")
```

### Product Monitoring

```go
// Schedule monitoring
result, err := client.ScheduleProductMonitoring("012345678901", "daily")

// Schedule multiple products
results, err := client.ScheduleProductMonitoringBatch(
    []string{"012345678901", "B08N5WRWNW"}, 
    "daily",
)

// Get scheduled products
scheduled, err := client.GetScheduledProducts()

// Remove from schedule
removed, err := client.RemoveProductFromSchedule("012345678901")

// Remove multiple from schedule
removedBatch, err := client.RemoveProductsFromSchedule([]string{"012345678901", "B08N5WRWNW"})
```

### Usage Information

```go
usage, err := client.GetUsage()
fmt.Printf("Credits remaining: %d\n", usage.Data.CreditsRemaining)
```

## Error Handling

The SDK provides specific error types for different scenarios:

```go
product, err := client.GetProductDetails("invalid-id")
if err != nil {
    switch e := err.(type) {
    case *shopsavvy.AuthenticationError:
        fmt.Println("Invalid API key")
    case *shopsavvy.NotFoundError:
        fmt.Println("Product not found")
    case *shopsavvy.ValidationError:
        fmt.Println("Invalid request parameters")
    case *shopsavvy.RateLimitError:
        fmt.Println("Rate limit exceeded")
    default:
        fmt.Printf("API error: %v\n", err)
    }
}
```

## Data Formats

All methods support optional format parameters for CSV output:

```go
// Get product details as CSV
product, err := client.GetProductDetails("012345678901", "csv")
```

## Response Structure

All API responses follow the same structure:

```go
type APIResponse[T any] struct {
    Success          bool   `json:"success"`
    Data             T      `json:"data"`
    Message          string `json:"message,omitempty"`
    CreditsUsed      *int   `json:"credits_used,omitempty"`
    CreditsRemaining *int   `json:"credits_remaining,omitempty"`
}
```

## Requirements

- Go 1.21 or higher
- Valid ShopSavvy API key (get one at [shopsavvy.com/data](https://shopsavvy.com/data))

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Links

- [ShopSavvy Data API](https://shopsavvy.com/data)
- [API Documentation](https://shopsavvy.com/data/documentation)
- [Get API Key](https://shopsavvy.com/data)
- [Report Issues](https://github.com/shopsavvy/sdk-go/issues)