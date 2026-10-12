# ShopSavvy Data API - Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/shopsavvy/sdk-go.svg)](https://pkg.go.dev/github.com/shopsavvy/sdk-go)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/docs-shopsavvy.com-blue)](https://shopsavvy.com/data/documentation)

Official Go SDK for the [ShopSavvy Data API](https://shopsavvy.com/data). Access comprehensive product data, real-time pricing, and historical price trends across **tens of thousands of retailers** and **30M+ products**.

## ⚡ 30-Second Quick Start

```go
// Install
// go get github.com/shopsavvy/sdk-go

// Use
package main

import (
    "fmt"
    "log"
    "github.com/shopsavvy/sdk-go"
)

func main() {
    client, _ := shopsavvy.NewClient("ss_live_your_api_key_here")
    defer client.Close()
    
    product, _ := client.GetProductDetails("012345678901")
    offers, _ := client.GetCurrentOffers("012345678901", "")
    
    bestPrice := offers.Data[0].Price
    bestRetailer := offers.Data[0].Retailer
    for _, offer := range offers.Data {
        if offer.Price < bestPrice {
            bestPrice = offer.Price
            bestRetailer = offer.Retailer
        }
    }
    
    fmt.Printf("%s - Best price: $%.2f at %s\n", product.Data.Name, bestPrice, bestRetailer)
}
```

## 🚀 Installation & Setup

### Installation

```bash
go get github.com/shopsavvy/sdk-go
```

### Get Your API Key

1. **Sign up**: Visit [shopsavvy.com/data](https://shopsavvy.com/data)
2. **Choose plan**: Select based on your usage needs  
3. **Get API key**: Copy from your dashboard
4. **Test**: Run the 30-second example above

## 📖 Complete API Reference

### Client Configuration

```go
import (
    "os"
    "time"
    "github.com/shopsavvy/sdk-go"
)

// Basic configuration
client, err := shopsavvy.NewClient("ss_live_your_api_key_here")

// Advanced configuration
client, err := shopsavvy.NewClient(
    "ss_live_your_api_key_here",
    shopsavvy.WithTimeout(60*time.Second),           // Request timeout
    shopsavvy.WithBaseURL("https://api.shopsavvy.com/v1"), // Custom base URL
    shopsavvy.WithRetryAttempts(3),                  // Retry failed requests
    shopsavvy.WithUserAgent("MyApp/1.0.0"),          // Custom user agent
)

// Environment variable configuration
apiKey := os.Getenv("SHOPSAVVY_API_KEY")
client, err := shopsavvy.NewClient(apiKey)
```

### Product Lookup

#### Single Product
```go
// Look up by barcode, ASIN, URL, model number, or ShopSavvy ID
product, err := client.GetProductDetails("012345678901")
amazonProduct, err := client.GetProductDetails("B08N5WRWNW")
urlProduct, err := client.GetProductDetails("https://www.amazon.com/dp/B08N5WRWNW")
modelProduct, err := client.GetProductDetails("MQ023LL/A") // iPhone model number

if err != nil {
    log.Fatal(err)
}

fmt.Printf("📦 Product: %s\n", product.Data.Name)
fmt.Printf("🏷️ Brand: %s\n", stringValue(product.Data.Brand))
fmt.Printf("📂 Category: %s\n", stringValue(product.Data.Category))
fmt.Printf("🔢 Product ID: %s\n", product.Data.ProductID)

// Helper function for optional strings
func stringValue(s *string) string {
    if s != nil {
        return *s
    }
    return "N/A"
}
```

#### Bulk Product Lookup
```go
// Up to 100 identifiers per request
identifiers := []string{
    "012345678901", "B08N5WRWNW", "045496590048",
    "https://www.bestbuy.com/site/product/123456",
    "MQ023LL/A", "SM-S911U", // iPhone and Samsung model numbers
}

products, err := client.GetProductDetailsBatch(identifiers)
if err != nil {
    log.Fatal(err)
}

for i, product := range products.Data {
    if product != nil {
        fmt.Printf("✓ Found: %s by %s\n", product.Name, stringValue(product.Brand))
    } else {
        fmt.Printf("❌ Failed to find product: %s\n", identifiers[i])
    }
}
```

### Real-Time Pricing

#### All Retailers Analysis

`GetCurrentOffers` returns one entry per product found; each carries its `Offers`. Offer
fields the API may omit (`Price`, `Retailer`, `Availability`, `Condition`, …) are pointers.
`Availability` is one of `"in"`, `"out"`, `"limited"`, `"pre-order"`, `"coming-soon"` or
`"discontinued"`, and nil when unknown.

```go
resp, err := client.GetCurrentOffers("012345678901", "")
if err != nil {
    log.Fatal(err)
}

for _, product := range resp.Data {
    // Keep only offers that have a price, cheapest first
    var priced []shopsavvy.Offer
    for _, offer := range product.Offers {
        if offer.Price != nil {
            priced = append(priced, offer)
        }
    }
    if len(priced) == 0 {
        continue
    }
    sort.Slice(priced, func(i, j int) bool { return *priced[i].Price < *priced[j].Price })

    name := func(o shopsavvy.Offer) string {
        if o.Retailer != nil {
            return *o.Retailer
        }
        return "unknown retailer"
    }
    cheapest, priciest := priced[0], priced[len(priced)-1]
    fmt.Printf("%s: %d offers\n", product.Title, len(product.Offers))
    fmt.Printf("💰 Best price: %s - $%.2f\n", name(cheapest), *cheapest.Price)
    fmt.Printf("💸 Highest price: %s - $%.2f\n", name(priciest), *priciest.Price)

    // Filter by availability and condition
    inStock, newCondition := 0, 0
    for _, offer := range product.Offers {
        if offer.Availability != nil && *offer.Availability == "in" {
            inStock++
        }
        if offer.Condition != nil && *offer.Condition == "new" {
            newCondition++
        }
    }
    fmt.Printf("✅ In-stock offers: %d\n", inStock)
    fmt.Printf("🆕 New condition: %d\n", newCondition)
}
```

#### Retailer-Specific Queries
```go
// Pass a retailer domain to get only that retailer's offers
retailers := []string{"amazon.com", "walmart.com", "target.com", "bestbuy.com"}

for _, retailer := range retailers {
    resp, err := client.GetCurrentOffers("012345678901", retailer)
    if err != nil {
        continue // Skip on error
    }
    for _, product := range resp.Data {
        for _, offer := range product.Offers {
            if offer.Price != nil {
                fmt.Printf("  %s: $%.2f\n", retailer, *offer.Price)
            }
        }
    }
}
```

### Price History

`GetPriceHistory` takes one or more identifiers (comma-separated), a `start` and `end`
date (`YYYY-MM-DD`, at most 366 days apart), and an optional retailer domain. `Data` holds
one entry per product found; each product carries its `Offers`, and each offer carries its
`History` points, newest first.

```go
history, err := client.GetPriceHistory("012345678901", "2024-01-01", "2024-01-31", "")
if err != nil {
    log.Fatal(err)
}

for _, product := range history.Data {
    fmt.Println(product.Title)
    for _, offer := range product.Offers {
        retailer := "unknown retailer"
        if offer.Retailer != nil {
            retailer = *offer.Retailer
        }
        fmt.Printf("  %s: %d price points\n", retailer, len(offer.History))
        for _, point := range offer.History {
            // Currency is "" on an archived point with no recorded currency;
            // Availability is "" when it was not known at that time.
            fmt.Printf("    %s  %.2f %s  %s\n", point.Timestamp, point.Price, point.Currency, point.Availability)
        }
    }
}

// Only one retailer's history:
amazonHistory, err := client.GetPriceHistory("012345678901", "2024-01-01", "2024-01-31", "amazon.com")
```

### Scheduled Refresh

Schedule products to be refreshed `"hourly"`, `"daily"` or `"weekly"`, optionally limited to
one retailer domain. Credit costs for each call are in the
[API documentation](https://shopsavvy.com/data/documentation).

```go
// Schedule one product (or several with ScheduleProductMonitoringBatch)
scheduled, err := client.ScheduleProductMonitoringBatch(
    []string{"012345678901", "B07GV2S1GS"}, "daily", "amazon.com")
if err != nil {
    log.Fatal(err)
}
// Data holds one ScheduledProductResult per product found: every product field
// (Title, ShopSavvy, Barcode, …) plus the Schedule that was set and the Retailer, if any.
for _, p := range scheduled.Data {
    fmt.Printf("%s (%s) → %s\n", p.Title, p.ShopSavvy, p.Schedule)
}
fmt.Printf("Credits used: %d\n", scheduled.CreditsUsed())

// List everything scheduled under this API key
list, err := client.GetScheduledProducts()
if err != nil {
    log.Fatal(err)
}
for _, p := range list.Data {
    // Schedule is nil for an interval with no Data API label (e.g. set via ShopSavvy
    // Business); Retailer is nil when refreshes cover every retailer.
    schedule, retailer := "custom", "all retailers"
    if p.Schedule != nil {
        schedule = *p.Schedule
    }
    if p.Retailer != nil {
        retailer = *p.Retailer
    }
    fmt.Printf("%s: %s, %s\n", p.Title, schedule, retailer)
}

// Stop refreshing (RemoveProductsFromSchedule for several). The API returns no data here.
removed, err := client.RemoveProductFromSchedule("012345678901")
if err != nil {
    log.Fatal(err)
}
fmt.Println(removed.Success, removed.Message) // true Products successfully removed from schedule
```

## 🚀 Production Deployment

### HTTP Server Integration

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/shopsavvy/sdk-go"
)

type PriceService struct {
    client *shopsavvy.Client
}

func NewPriceService() *PriceService {
    client, err := shopsavvy.NewClient(os.Getenv("SHOPSAVVY_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }
    return &PriceService{client: client}
}

func (ps *PriceService) GetOffersHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    identifier := vars["identifier"]
    
    offers, err := ps.client.GetCurrentOffers(identifier, "")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Find best price
    var bestPrice float64 = 999999
    for _, offer := range offers.Data {
        if offer.Price < bestPrice {
            bestPrice = offer.Price
        }
    }
    
    response := map[string]interface{}{
        "success":           true,
        "product_id":        identifier,
        "offers":           offers.Data,
        "best_price":       bestPrice,
        "credits_remaining": offers.CreditsRemaining,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    service := NewPriceService()
    defer service.client.Close()
    
    r := mux.NewRouter()
    r.HandleFunc("/api/products/{identifier}/offers", service.GetOffersHandler).Methods("GET")
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

### Background Worker

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/shopsavvy/sdk-go"
)

type PriceMonitor struct {
    client   *shopsavvy.Client
    products []string
}

func NewPriceMonitor() *PriceMonitor {
    client, err := shopsavvy.NewClient(os.Getenv("SHOPSAVVY_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }
    
    return &PriceMonitor{
        client:   client,
        products: []string{"012345678901", "B08N5WRWNW"},
    }
}

func (pm *PriceMonitor) Start(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            pm.checkPrices()
        }
    }
}

func (pm *PriceMonitor) checkPrices() {
    for _, productID := range pm.products {
        offers, err := pm.client.GetCurrentOffers(productID, "")
        if err != nil {
            log.Printf("Error checking prices for %s: %v", productID, err)
            continue
        }
        
        if len(offers.Data) > 0 {
            bestPrice := offers.Data[0].Price
            bestRetailer := offers.Data[0].Retailer
            
            for _, offer := range offers.Data {
                if offer.Price < bestPrice {
                    bestPrice = offer.Price
                    bestRetailer = offer.Retailer
                }
            }
            
            log.Printf("Product %s: Best price $%.2f at %s", productID, bestPrice, bestRetailer)
            // Here you could send alerts, update database, etc.
        }
    }
}

func main() {
    monitor := NewPriceMonitor()
    defer monitor.client.Close()
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    log.Println("Starting price monitor...")
    monitor.Start(ctx)
}
```

## 🛠️ Development & Testing

### Local Development Setup

```bash
# Clone the repository
git clone https://github.com/shopsavvy/sdk-go.git
cd sdk-go

# Download dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build

# Run examples
go run examples/basic_usage.go
```

### Testing Your Integration

```go
package main

import (
    "fmt"
    "log"
    "github.com/shopsavvy/sdk-go"
)

func main() {
    // Use test API key (starts with ss_test_)
    client, err := shopsavvy.NewClient("ss_test_your_test_key_here")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Test product lookup
    product, err := client.GetProductDetails("012345678901")
    if err != nil {
        fmt.Printf("❌ Test failed: %v\n", err)
        return
    }
    fmt.Printf("✅ Product lookup: %s\n", product.Data.Name)
    
    // Test current offers
    offers, err := client.GetCurrentOffers("012345678901", "")
    if err != nil {
        fmt.Printf("❌ Test failed: %v\n", err)
        return
    }
    fmt.Printf("✅ Current offers: %d found\n", len(offers.Data))
    
    // Test usage info
    usage, err := client.GetUsage()
    if err != nil {
        fmt.Printf("❌ Test failed: %v\n", err)
        return
    }
    fmt.Printf("✅ API usage: %d credits remaining\n", *usage.Data.CreditsRemaining)
    
    fmt.Println("\n🎉 All tests passed! SDK is working correctly.")
}
```

## 📚 Additional Resources

- **[ShopSavvy Data API Documentation](https://shopsavvy.com/data/documentation)** - Complete API reference
- **[API Dashboard](https://shopsavvy.com/data/dashboard)** - Manage your API keys and usage
- **[GitHub Repository](https://github.com/shopsavvy/sdk-go)** - Source code and issues
- **[Go Package Documentation](https://pkg.go.dev/github.com/shopsavvy/sdk-go)** - Package reference
- **[Support](mailto:business@shopsavvy.com)** - Get help from our team

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## About ShopSavvy

ShopSavvy has been comparing prices since 2008, and the app has 40M+ downloads. The ShopSavvy Data API serves the same product and price data over REST: send a barcode (UPC/EAN/ISBN), ASIN, product URL, model number or MPN, and get back the product record, current offers from tens of thousands of retailers, and price history back to 2011 (dense from 2019).

- **Lookup by any identifier**: barcode, ASIN, URL, model number, MPN or ShopSavvy ID, up to 100 per request
- **Current offers**: price, availability, condition, seller and when each offer was last seen
- **Price history**: quality-filtered price points per retailer over any date range
- **Scheduled refreshes**: hourly, daily or weekly re-checks for the products you care about, at retailers that can be re-checked through an official retailer API (others update as ShopSavvy sees them)
- **JSON or CSV** responses, with credits used and remaining on every response
- **Open-source SDKs in 10 languages**: Python, TypeScript, Go, Rust, Java, Kotlin, Swift, PHP, Ruby and C#

Plans and the 14-day trial are at [shopsavvy.com/data](https://shopsavvy.com/data); full reference at [shopsavvy.com/data/documentation](https://shopsavvy.com/data/documentation).

---

**Ready to get started?** [Sign up for your API key](https://shopsavvy.com/data) • **Need help?** [Contact us](mailto:business@shopsavvy.com)