package shopsavvy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// newFixtureServer serves testdata/<fixture> for requests to wantPath and records the
// query string of every request it receives.
func newFixtureServer(t *testing.T, wantPath, fixture string, gotQuery *map[string]string) *httptest.Server {
	t.Helper()
	body, err := os.ReadFile("testdata/" + fixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, wantPath)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got, want := r.Header.Get("User-Agent"), "ShopSavvy-Go-SDK/"+Version; got != want {
			t.Errorf("User-Agent = %q, want %q", got, want)
		}
		q := map[string]string{}
		for k, v := range r.URL.Query() {
			q[k] = v[0]
		}
		*gotQuery = q
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	c, err := NewClient("ss_test_abc123", WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func strOrNil(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}

func TestGetPriceHistorySendsStartEnd(t *testing.T) {
	var q map[string]string
	srv := newFixtureServer(t, "/products/offers/history", "price-history-fixture.json", &q)
	c := newTestClient(t, srv.URL)

	if _, err := c.GetPriceHistory("611247373064,611247369449", "2022-11-20", "2022-11-27", "amazon.com"); err != nil {
		t.Fatalf("GetPriceHistory: %v", err)
	}
	want := map[string]string{"ids": "611247373064,611247369449", "start": "2022-11-20", "end": "2022-11-27", "retailer": "amazon.com"}
	for k, v := range want {
		if q[k] != v {
			t.Errorf("query %s = %q, want %q", k, q[k], v)
		}
	}
	for _, legacy := range []string{"start_date", "end_date"} {
		if _, ok := q[legacy]; ok {
			t.Errorf("legacy query param %q must not be sent", legacy)
		}
	}
	if _, ok := q["format"]; ok {
		t.Errorf("format must not be sent when not requested")
	}
}

func TestGetPriceHistoryOmitsEmptyRetailer(t *testing.T) {
	var q map[string]string
	srv := newFixtureServer(t, "/products/offers/history", "price-history-fixture.json", &q)
	c := newTestClient(t, srv.URL)

	if _, err := c.GetPriceHistory("611247373064", "2022-11-20", "2022-11-27", ""); err != nil {
		t.Fatalf("GetPriceHistory: %v", err)
	}
	if _, ok := q["retailer"]; ok {
		t.Errorf("retailer must be omitted when empty, got %q", q["retailer"])
	}
}

func TestGetPriceHistoryParsesProductsOffersAndHistory(t *testing.T) {
	var q map[string]string
	srv := newFixtureServer(t, "/products/offers/history", "price-history-fixture.json", &q)
	c := newTestClient(t, srv.URL)

	resp, err := c.GetPriceHistory("611247373064,611247369449", "2022-11-20", "2022-11-27", "")
	if err != nil {
		t.Fatalf("GetPriceHistory: %v", err)
	}

	if !resp.Success {
		t.Errorf("Success = false")
	}
	if resp.CreditsUsed() != 14 || resp.CreditsRemaining() != 986 {
		t.Errorf("credits used/remaining = %d/%d, want 14/986", resp.CreditsUsed(), resp.CreditsRemaining())
	}
	if resp.Meta == nil || resp.Meta.RateLimitRemaining == nil || *resp.Meta.RateLimitRemaining != 999 {
		t.Errorf("rate_limit_remaining not parsed: %+v", resp.Meta)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("len(Data) = %d, want 2 products", len(resp.Data))
	}

	// Product 1: full product fields + two offers with history.
	p1 := resp.Data[0]
	if p1.Title != "Keurig K-Mini Single Serve Coffee Maker, Black" || p1.ShopSavvy != "3ONn300xybP3y66ibqc1" {
		t.Errorf("product 1 title/shopsavvy = %q/%q", p1.Title, p1.ShopSavvy)
	}
	if strOrNil(p1.Barcode) != "611247373064" || strOrNil(p1.Amazon) != "B07G14HTBZ" || strOrNil(p1.Brand) != "Keurig" {
		t.Errorf("product 1 barcode/amazon/brand = %s/%s/%s", strOrNil(p1.Barcode), strOrNil(p1.Amazon), strOrNil(p1.Brand))
	}
	if strOrNil(p1.TitleShort) != "Keurig K-Mini" || p1.Rating == nil || p1.Rating.Count != 51234 {
		t.Errorf("product 1 title_short/rating = %s/%+v", strOrNil(p1.TitleShort), p1.Rating)
	}
	if len(p1.Images) != 1 || p1.Identifiers["upc"] != "611247373064" {
		t.Errorf("product 1 images/identifiers = %v/%v", p1.Images, p1.Identifiers)
	}
	if len(p1.Offers) != 2 {
		t.Fatalf("product 1 offers = %d, want 2", len(p1.Offers))
	}

	amazon := p1.Offers[0]
	if amazon.ID != "0IUouCFtZEhxeOablTPl" || strOrNil(amazon.Retailer) != "Amazon" || strOrNil(amazon.Seller) != "ACME Deals" {
		t.Errorf("amazon offer id/retailer/seller = %q/%s/%s", amazon.ID, strOrNil(amazon.Retailer), strOrNil(amazon.Seller))
	}
	if amazon.Price == nil || *amazon.Price != 74.96 || strOrNil(amazon.Currency) != "USD" || strOrNil(amazon.Availability) != "in" {
		t.Errorf("amazon offer price/currency/availability wrong: %+v", amazon.Offer)
	}
	if strOrNil(amazon.URL) != "https://www.amazon.com/dp/B07G14HTBZ?m=A1GKQADQC2VI6E" || strOrNil(amazon.Timestamp) != "2022-11-27T22:36:33.236Z" {
		t.Errorf("amazon offer URL/timestamp = %s/%s", strOrNil(amazon.URL), strOrNil(amazon.Timestamp))
	}
	if len(amazon.History) != 3 {
		t.Fatalf("amazon history = %d points, want 3", len(amazon.History))
	}
	wantHistory := []PriceHistoryEntry{
		{Timestamp: "2022-11-27T22:36:33.236Z", Price: 74.96, Currency: "USD", Availability: "in"},
		{Timestamp: "2022-11-24T10:00:00.000Z", Price: 70.99, Currency: "USD", Availability: "out"},
		// JSON null currency and missing availability both decode to "".
		{Timestamp: "2022-11-21T08:15:00.000Z", Price: 79.99, Currency: "", Availability: ""},
	}
	for i, want := range wantHistory {
		if amazon.History[i] != want {
			t.Errorf("amazon history[%d] = %+v, want %+v", i, amazon.History[i], want)
		}
	}

	bestBuy := p1.Offers[1]
	if bestBuy.ID != "Z9kQ2mBestBuyOffer01" || strOrNil(bestBuy.Retailer) != "Best Buy" {
		t.Errorf("best buy offer id/retailer = %q/%s", bestBuy.ID, strOrNil(bestBuy.Retailer))
	}
	if bestBuy.Seller != nil {
		t.Errorf("best buy seller = %q, want nil (JSON null)", *bestBuy.Seller)
	}
	if bestBuy.Availability != nil {
		t.Errorf("best buy availability = %q, want nil (key absent)", *bestBuy.Availability)
	}
	if len(bestBuy.History) != 2 || bestBuy.History[1].Price != 64.99 || bestBuy.History[1].Timestamp != "2022-11-22T11:47:00.226Z" {
		t.Errorf("best buy history = %+v", bestBuy.History)
	}

	// Product 2: null optional product fields, one offer with an empty history.
	p2 := resp.Data[1]
	if p2.ShopSavvy != "DrKWneG0MpFlZpwZXNYa" || strOrNil(p2.Model) != "K-Elite" {
		t.Errorf("product 2 shopsavvy/model = %q/%s", p2.ShopSavvy, strOrNil(p2.Model))
	}
	if p2.Category != nil || p2.Amazon != nil || p2.Color != nil || p2.MPN != nil {
		t.Errorf("product 2 null fields must decode to nil: category=%s amazon=%s color=%s mpn=%s",
			strOrNil(p2.Category), strOrNil(p2.Amazon), strOrNil(p2.Color), strOrNil(p2.MPN))
	}
	if p2.Rating != nil || len(p2.Images) != 0 {
		t.Errorf("product 2 rating/images = %+v/%v", p2.Rating, p2.Images)
	}
	if len(p2.Offers) != 1 {
		t.Fatalf("product 2 offers = %d, want 1", len(p2.Offers))
	}
	ebay := p2.Offers[0]
	if ebay.ID != "eBayListing000000001" || strOrNil(ebay.Condition) != "used" || ebay.Price == nil || *ebay.Price != 89.5 {
		t.Errorf("ebay offer = %+v", ebay.Offer)
	}
	if ebay.History == nil || len(ebay.History) != 0 {
		t.Errorf("ebay history = %#v, want empty non-nil slice", ebay.History)
	}
}

type recordedRequest struct {
	method string
	path   string
	query  string
	body   string
}

// newRecordingServer answers every request with body and records what the client sent.
func newRecordingServer(t *testing.T, body string, got *recordedRequest) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		*got = recordedRequest{method: r.Method, path: r.URL.Path, query: r.URL.RawQuery, body: string(b)}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// scheduleResponseBody is the shape the schedule handler returns: one entry per product,
// the product's public fields plus `schedule` (and `retailer` when one was requested).
const scheduleResponseBody = `{"success":true,"data":[
 {"title":"Keurig K-Mini","shopsavvy":"3ONn300xybP3y66ibqc1","barcode":"611247373064","amazon":null,"images":[],"schedule":"daily","retailer":"amazon.com"},
 {"title":"Keurig K-Elite","shopsavvy":"DrKWneG0MpFlZpwZXNYa","barcode":"611247369449","images":[],"schedule":"daily","retailer":"amazon.com"}],
 "meta":{"request_id":"r1","credits_used":2,"credits_remaining":998,"rate_limit_remaining":999}}`

const unscheduleResponseBody = `{"success":true,"message":"Products successfully removed from schedule","meta":{"request_id":"r2","credits_used":0,"credits_remaining":0,"rate_limit_remaining":0}}`

func assertRequest(t *testing.T, got recordedRequest, method, path string, query url.Values) {
	t.Helper()
	if got.method != method || got.path != path {
		t.Errorf("request = %s %s, want %s %s", got.method, got.path, method, path)
	}
	gotQuery, err := url.ParseQuery(got.query)
	if err != nil {
		t.Fatalf("parse query %q: %v", got.query, err)
	}
	if gotQuery.Encode() != query.Encode() {
		t.Errorf("query = %q, want %q", gotQuery.Encode(), query.Encode())
	}
	if got.body != "" {
		t.Errorf("request body = %q, want none (the server reads only query params)", got.body)
	}
}

func TestScheduleProductMonitoringSendsQueryParams(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, scheduleResponseBody, &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.ScheduleProductMonitoring("611247373064", "daily")
	if err != nil {
		t.Fatalf("ScheduleProductMonitoring: %v", err)
	}
	assertRequest(t, got, http.MethodPut, "/products/scheduled", url.Values{"ids": {"611247373064"}, "schedule": {"daily"}})
	if !strings.Contains(got.query, "ids=611247373064") {
		t.Errorf("raw query %q", got.query)
	}

	if len(resp.Data) != 2 || resp.Data[0].ShopSavvy != "3ONn300xybP3y66ibqc1" || resp.Data[0].Schedule != "daily" {
		t.Fatalf("parsed data = %+v", resp.Data)
	}
	if strOrNil(resp.Data[0].Retailer) != "amazon.com" || resp.Data[0].Amazon != nil || strOrNil(resp.Data[1].Barcode) != "611247369449" {
		t.Errorf("parsed product fields wrong: %+v", resp.Data)
	}
	if resp.CreditsUsed() != 2 {
		t.Errorf("CreditsUsed = %d, want 2", resp.CreditsUsed())
	}
}

func TestScheduleProductMonitoringWithRetailer(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, scheduleResponseBody, &got)
	c := newTestClient(t, srv.URL)

	if _, err := c.ScheduleProductMonitoring("611247373064", "hourly", "amazon.com"); err != nil {
		t.Fatalf("ScheduleProductMonitoring: %v", err)
	}
	assertRequest(t, got, http.MethodPut, "/products/scheduled", url.Values{"ids": {"611247373064"}, "schedule": {"hourly"}, "retailer": {"amazon.com"}})
}

func TestScheduleProductMonitoringBatchSendsQueryParams(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, scheduleResponseBody, &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.ScheduleProductMonitoringBatch([]string{"611247373064", "611247369449"}, "weekly")
	if err != nil {
		t.Fatalf("ScheduleProductMonitoringBatch: %v", err)
	}
	assertRequest(t, got, http.MethodPut, "/products/scheduled", url.Values{"ids": {"611247373064,611247369449"}, "schedule": {"weekly"}})
	if !strings.Contains(got.query, "ids=611247373064%2C611247369449") {
		t.Errorf("ids not comma-joined and URL-encoded in raw query %q", got.query)
	}
	if len(resp.Data) != 2 || resp.Data[1].Title != "Keurig K-Elite" {
		t.Errorf("parsed data = %+v", resp.Data)
	}
}

func TestScheduleProductMonitoringBatchWithRetailer(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, scheduleResponseBody, &got)
	c := newTestClient(t, srv.URL)

	if _, err := c.ScheduleProductMonitoringBatch([]string{"611247373064", "B07G14HTBZ"}, "daily", "bestbuy.com"); err != nil {
		t.Fatalf("ScheduleProductMonitoringBatch: %v", err)
	}
	assertRequest(t, got, http.MethodPut, "/products/scheduled", url.Values{"ids": {"611247373064,B07G14HTBZ"}, "schedule": {"daily"}, "retailer": {"bestbuy.com"}})
}

func TestRemoveProductFromScheduleSendsQueryParams(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, unscheduleResponseBody, &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.RemoveProductFromSchedule("611247373064")
	if err != nil {
		t.Fatalf("RemoveProductFromSchedule: %v", err)
	}
	assertRequest(t, got, http.MethodDelete, "/products/scheduled", url.Values{"ids": {"611247373064"}})
	if !resp.Success || resp.Message != "Products successfully removed from schedule" {
		t.Errorf("parsed response = %+v", resp)
	}
}

func TestRemoveProductsFromScheduleSendsQueryParams(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, unscheduleResponseBody, &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.RemoveProductsFromSchedule([]string{"611247373064", "611247369449"})
	if err != nil {
		t.Fatalf("RemoveProductsFromSchedule: %v", err)
	}
	assertRequest(t, got, http.MethodDelete, "/products/scheduled", url.Values{"ids": {"611247373064,611247369449"}})
	if !resp.Success {
		t.Errorf("Success = false")
	}
}

func TestGetCurrentOffersParsesInStorePrices(t *testing.T) {
	var q map[string]string
	srv := newFixtureServer(t, "/products/offers", "offers-in-store-fixture.json", &q)
	c := newTestClient(t, srv.URL)

	resp, err := c.GetCurrentOffers("011110838360", "")
	if err != nil {
		t.Fatalf("GetCurrentOffers: %v", err)
	}
	if q["ids"] != "011110838360" {
		t.Errorf("ids = %q", q["ids"])
	}
	product := resp.Data[0]
	if len(product.Offers) != 1 || len(product.InStorePrices) != 1 {
		t.Fatalf("offers/in-store = %d/%d, want 1/1", len(product.Offers), len(product.InStorePrices))
	}
	inStore := product.InStorePrices[0]
	if inStore.Source != "user-collected" || inStore.Retailer == nil || *inStore.Retailer != "Kroger" {
		t.Errorf("source/retailer = %q/%v", inStore.Source, inStore.Retailer)
	}
	if inStore.Store.Address == nil || *inStore.Store.Address != "1014 Vine St, Cincinnati, OH 45202, USA" {
		t.Errorf("address = %v", inStore.Store.Address)
	}
	if inStore.Store.Latitude == nil || *inStore.Store.Latitude != 39.1086 {
		t.Errorf("latitude = %v", inStore.Store.Latitude)
	}
	if inStore.Price != 3.49 || inStore.Currency != "USD" || inStore.ObservedOn == nil || *inStore.ObservedOn != "2026-11-07" {
		t.Errorf("price/currency/observed_on = %v/%q/%v", inStore.Price, inStore.Currency, inStore.ObservedOn)
	}
}
