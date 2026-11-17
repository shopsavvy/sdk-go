# Changelog

All notable changes to this project are documented in this file. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## 1.6.0 - 2026-11-17

### Added
- Every API error (`APIError`, `AuthenticationError`, `NotFoundError`, `ValidationError`, `RateLimitError`) now embeds `ErrorDetails` with `Code`, `Action`, `DocsURL`, `RequestID` and `RetryAfterSeconds` (parsed from the `Retry-After` header), read from the API's `error_detail` and `meta.request_id`.
- `CreditsExhaustedError`, returned for `ERR_CREDITS_EXHAUSTED`. It is a separate type from `RateLimitError`.
- README section "Errors, rate limits and credits".

### Changed
- Behavioral change: an HTTP 429 with `ERR_CREDITS_EXHAUSTED` is no longer reported as `RateLimitError`; it is a `CreditsExhaustedError`. Code that retried on `RateLimitError` no longer retries when credits are used up. A 429 without a code still returns `RateLimitError`.
- README quick start reads the key from `SHOPSAVVY_API_KEY` and prints the cheapest offer for a real barcode; removed the "sign up for a free account" and "most comprehensive" wording.

## 1.5.0 - 2026-11-09

### Added
- `in_store_prices` (user-collected receipt prices) on current offers and batch results.

## 1.4.0 - 2026-09-27

### Added
- Currency on price-history points.

### Fixed
- `GetPriceHistory` decodes the products with nested offers response, and sends `start`/`end`.
- Schedule and unschedule send query parameters to `/products/scheduled`, and model the real list and unschedule responses.

## 1.1.0 - 2026-03-31

### Added
- Deals and reviews endpoints and expanded product types.
