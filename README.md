# ShopSavvy Data API - Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/shopsavvy/sdk-go.svg)](https://pkg.go.dev/github.com/shopsavvy/sdk-go)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/docs-shopsavvy.com-blue)](https://shopsavvy.com/data/documentation)

Official Go SDK for the [ShopSavvy Data API](https://shopsavvy.com/data). Access comprehensive product data, real-time pricing, and historical price trends across **thousands of retailers** and **millions of products**.

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

## 📊 Feature Comparison

| Feature | Free Tier | Pro | Enterprise |
|---------|-----------|-----|-----------| 
| **API Calls/Month** | 1,000 | 100,000 | Unlimited |
| **Product Details** | ✅ | ✅ | ✅ |
| **Real-time Pricing** | ✅ | ✅ | ✅ |
| **Price History** | 30 days | 1 year | 5+ years |
| **Bulk Operations** | 10/batch | 100/batch | 1000/batch |
| **Retailer Coverage** | 50+ | 500+ | 1000+ |
| **Rate Limiting** | 60/hour | 1000/hour | Custom |
| **Support** | Community | Email | Phone + Dedicated |

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
// Process up to 100 products at once (Pro plan)
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
```go
offers, err := client.GetCurrentOffers("012345678901", "")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d offers across retailers\n", len(offers.Data))

// Advanced price analysis
sort.Slice(offers.Data, func(i, j int) bool {
    return offers.Data[i].Price < offers.Data[j].Price
})

cheapest := offers.Data[0]
mostExpensive := offers.Data[len(offers.Data)-1]

var total float64
for _, offer := range offers.Data {
    total += offer.Price
}
average := total / float64(len(offers.Data))

fmt.Printf("💰 Best price: %s - $%.2f\n", cheapest.Retailer, cheapest.Price)
fmt.Printf("💸 Highest price: %s - $%.2f\n", mostExpensive.Retailer, mostExpensive.Price)
fmt.Printf("📊 Average price: $%.2f\n", average)
fmt.Printf("💡 Potential savings: $%.2f\n", mostExpensive.Price-cheapest.Price)

// Filter by availability and condition
inStockCount := 0
newConditionCount := 0
for _, offer := range offers.Data {
    if offer.Availability == "in_stock" {
        inStockCount++
    }
    if offer.Condition == "new" {
        newConditionCount++
    }
}

fmt.Printf("✅ In-stock offers: %d\n", inStockCount)
fmt.Printf("🆕 New condition: %d\n", newConditionCount)
```

#### Retailer-Specific Queries
```go
// Major retailers
retailers := []string{"amazon", "walmart", "target", "bestbuy"}
retailerPrices := make(map[string]float64)

for _, retailer := range retailers {
    offers, err := client.GetCurrentOffers("012345678901", retailer)
    if err != nil {
        continue // Skip on error
    }
    
    if len(offers.Data) > 0 {
        // Find best offer for this retailer
        bestPrice := offers.Data[0].Price
        for _, offer := range offers.Data {
            if offer.Price < bestPrice {
                bestPrice = offer.Price
            }
        }
        retailerPrices[retailer] = bestPrice
    }
}

// Sort and display results
type retailerPrice struct {
    name  string
    price float64
}

var sorted []retailerPrice
for name, price := range retailerPrices {
    sorted = append(sorted, retailerPrice{name, price})
}

sort.Slice(sorted, func(i, j int) bool {
    return sorted[i].price < sorted[j].price
})

fmt.Println("Retailer price comparison:")
for _, rp := range sorted {
    fmt.Printf("  %s: $%.2f\n", strings.Title(rp.name), rp.price)
}
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

## 🏢 About ShopSavvy

**ShopSavvy** is the world's first mobile shopping app, helping consumers find the best deals since 2008. With over **40 million downloads** and millions of active users, ShopSavvy has saved consumers billions of dollars.

### Why Choose ShopSavvy Data API?
- ✅ **Trusted by millions** - Proven at scale since 2008
- ✅ **Comprehensive coverage** - 1000+ retailers, millions of products  
- ✅ **Real-time accuracy** - Fresh data updated continuously
- ✅ **Developer-friendly** - Easy integration, great documentation
- ✅ **Reliable infrastructure** - 99.9% uptime, enterprise-grade
- ✅ **Flexible pricing** - Plans for every use case and budget

---

**Ready to get started?** [Sign up for your API key](https://shopsavvy.com/data) • **Need help?** [Contact us](mailto:business@shopsavvy.com)