# Go FT HTTP Client

Utility package for FT specific http client code.

# FTHttp

The fthttp package provides functionality to create new `http.Client` in a standardized way.
Supported configuration options:
* Timeout - value will be taken in to account as a whole which would also cover the time spent reading the response body.
* Logging - will create a log entry for each request 
* UserAgent - will add the correct user agent to the outgoing request.

# Transport

The transport package contains an `http.RoundTripper` implementation which allows simple modifications to `http.Request` via extensions before delegating to the core lib `http.DefaultRoundTripper`.
It also allows the RoundTripper to be instrumented to start logging the request it makes.

There are currently two extensions currently implemented:

* `WithStandardUserAgent` which uses provided platform and system code values to append an RFC7231 compliant `User-Agent` header in the format: `PLATFORM-system-code/version`
* `TransactionIDFromContext` will attempt to retrieve the transaction ID from the `*http.Request` context, and if it finds one will set the `X-Request-Id` header.

**Important** Neither extension will override an existing header, it will only add a new header if there is none.

Initialize logging with `WithLogger` to start logging the outgoing requests.

# Usage

## HTTP Client

For examples how to create new fthttp client refer to `examples_test.go`

## Transport

For examples how to use `transport.ExtensibleTransport` refer to the examples in `examples_test.go`

# Notable Dependencies

*  [go-logger v2](https://github.com/Financial-Times/go-logger)
