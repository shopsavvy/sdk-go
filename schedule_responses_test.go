package shopsavvy

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

// These tests drive the real client methods against full-size fixtures of the exact
// JSON the schedule / unschedule / scheduled handlers in the Data API emit
// (publicProductForProduct fields + schedule/retailer, and the meta block).

func fixtureBody(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(b)
}

func TestScheduleParsesFullProductFieldsAndMeta(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, fixtureBody(t, "schedule-put-fixture.json"), &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.ScheduleProductMonitoringBatch([]string{"611247373064", "611247369449"}, "daily", "amazon.com")
	if err != nil {
		t.Fatalf("ScheduleProductMonitoringBatch: %v", err)
	}
	assertRequest(t, got, http.MethodPut, "/products/scheduled", url.Values{"ids": {"611247373064,611247369449"}, "schedule": {"daily"}, "retailer": {"amazon.com"}})

	if !resp.Success || len(resp.Data) != 2 {
		t.Fatalf("success/len = %v/%d", resp.Success, len(resp.Data))
	}
	p := resp.Data[0]
	if p.Title != "Keurig K-Mini Single Serve Coffee Maker, Black" || p.ShopSavvy != "3ONn300xybP3y66ibqc1" {
		t.Errorf("title/shopsavvy = %q/%q", p.Title, p.ShopSavvy)
	}
	if p.Schedule != "daily" || strOrNil(p.Retailer) != "amazon.com" {
		t.Errorf("schedule/retailer = %q/%s", p.Schedule, strOrNil(p.Retailer))
	}
	if strOrNil(p.Amazon) != "B07GV2S1GS" || strOrNil(p.Barcode) != "611247373064" || strOrNil(p.MPN) != "5000200237" || strOrNil(p.Model) != "K-Mini" {
		t.Errorf("identifiers wrong: %+v", p.ProductDetails)
	}
	if len(p.Images) != 2 || p.Images[1] != "https://images.shopsavvy.com/3ONn300xybP3y66ibqc1/side.jpg" {
		t.Errorf("images = %v", p.Images)
	}
	if strOrNil(p.TitleShort) != "Keurig K-Mini" || strOrNil(p.Slug) != "keurig-k-mini-single-serve-coffee-maker-black" || strOrNil(p.Description) != "Brews any cup size between 6 and 12 oz." {
		t.Errorf("expanded text fields wrong: %+v", p.ProductDetails)
	}
	if len(p.Categories) != 2 || p.Attributes["Wattage"] != "600 W" || len(p.Keywords) != 2 {
		t.Errorf("categories/attributes/keywords = %v/%v/%v", p.Categories, p.Attributes, p.Keywords)
	}
	if p.Rating == nil || p.Rating.Value != 4.6 || p.Rating.Count != 18234 {
		t.Errorf("rating = %+v", p.Rating)
	}
	if p.Score == nil || p.Score.Overall == nil || *p.Score.Overall != 0.82 || p.Score.Aspects["reliability"] != 0.7 {
		t.Errorf("score = %+v", p.Score)
	}
	if p.Identifiers["walmart"] != "592453934" {
		t.Errorf("identifiers = %v", p.Identifiers)
	}

	p2 := resp.Data[1]
	if p2.Category != nil || p2.Color != nil || p2.MPN != nil || p2.Amazon != nil || p2.Rating != nil || p2.Score != nil {
		t.Errorf("null/absent product fields should be nil: %+v", p2.ProductDetails)
	}
	if p2.Schedule != "daily" || strOrNil(p2.Retailer) != "amazon.com" || len(p2.Images) != 0 {
		t.Errorf("product 2 schedule/retailer/images = %q/%s/%v", p2.Schedule, strOrNil(p2.Retailer), p2.Images)
	}

	if resp.Meta == nil || strOrNil(resp.Meta.RequestID) != "req_put_7f3a" || resp.CreditsUsed() != 2 || resp.CreditsRemaining() != 998 {
		t.Errorf("meta = %+v", resp.Meta)
	}
	if resp.Meta.RateLimitRemaining == nil || *resp.Meta.RateLimitRemaining != 59 {
		t.Errorf("rate_limit_remaining = %v", resp.Meta.RateLimitRemaining)
	}
}

func TestGetScheduledProductsParsesProductFieldsScheduleAndRetailer(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, fixtureBody(t, "schedule-list-fixture.json"), &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.GetScheduledProducts()
	if err != nil {
		t.Fatalf("GetScheduledProducts: %v", err)
	}
	if got.method != http.MethodGet || got.path != "/products/scheduled" || got.query != "" {
		t.Errorf("request = %s %s?%s", got.method, got.path, got.query)
	}
	if !resp.Success || len(resp.Data) != 3 {
		t.Fatalf("success/len = %v/%d", resp.Success, len(resp.Data))
	}

	first := resp.Data[0]
	if first.ShopSavvy != "3ONn300xybP3y66ibqc1" || first.Title != "Keurig K-Mini Single Serve Coffee Maker, Black" || strOrNil(first.Barcode) != "611247373064" {
		t.Errorf("first product = %+v", first.ProductDetails)
	}
	if strOrNil(first.Schedule) != "hourly" || strOrNil(first.Retailer) != "bestbuy.com" {
		t.Errorf("first schedule/retailer = %s/%s", strOrNil(first.Schedule), strOrNil(first.Retailer))
	}
	if first.Rating == nil || first.Rating.Count != 18234 || first.Identifiers["amazon"] != "B07GV2S1GS" {
		t.Errorf("first rating/identifiers = %+v/%v", first.Rating, first.Identifiers)
	}

	// Scheduled for every retailer: no `retailer` key.
	second := resp.Data[1]
	if strOrNil(second.Schedule) != "weekly" || second.Retailer != nil || second.Category != nil || strOrNil(second.Model) != "K-Elite" {
		t.Errorf("second = %+v schedule=%s retailer=%v", second.ProductDetails, strOrNil(second.Schedule), second.Retailer)
	}

	// Saved with an interval that has no Data API label (e.g. 4h from ShopSavvy Business):
	// the server omits `schedule` entirely.
	third := resp.Data[2]
	if third.Schedule != nil || third.Retailer != nil || third.ShopSavvy != "Qm1bXz0AirPodsPro02" {
		t.Errorf("third = %+v schedule=%v retailer=%v", third.ProductDetails, third.Schedule, third.Retailer)
	}

	if resp.Meta == nil || strOrNil(resp.Meta.RequestID) != "req_list_91bc" || resp.CreditsUsed() != 0 {
		t.Errorf("meta = %+v", resp.Meta)
	}
}

func TestRemoveFromScheduleParsesSuccessMessageAndMeta(t *testing.T) {
	var got recordedRequest
	srv := newRecordingServer(t, fixtureBody(t, "schedule-delete-fixture.json"), &got)
	c := newTestClient(t, srv.URL)

	resp, err := c.RemoveProductsFromSchedule([]string{"611247373064", "611247369449"})
	if err != nil {
		t.Fatalf("RemoveProductsFromSchedule: %v", err)
	}
	assertRequest(t, got, http.MethodDelete, "/products/scheduled", url.Values{"ids": {"611247373064,611247369449"}})
	if !resp.Success || resp.Message != "Products successfully removed from schedule" {
		t.Errorf("success/message = %v/%q", resp.Success, resp.Message)
	}
	if resp.Meta == nil || strOrNil(resp.Meta.RequestID) != "req_del_22d0" || resp.Meta.CreditsUsed != 0 || resp.Meta.CreditsRemaining != 0 {
		t.Errorf("meta = %+v", resp.Meta)
	}
	if resp.Meta.RateLimitRemaining == nil || *resp.Meta.RateLimitRemaining != 0 {
		t.Errorf("rate_limit_remaining = %v", resp.Meta.RateLimitRemaining)
	}
}
