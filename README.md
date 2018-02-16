# Go FT HTTP Client

Utility package for FT specific http client code.

# Transport

The transport package contains an `http.RoundTripper` implementation which allows simple modifications to `http.Request` via extensions before delegating to the core lib `http.DefaultRoundTripper`.

There are currently two extensions currently implemented:

* `WithStandardUserAgent` which uses provided platform and system code values to append an RFC7231 compliant `User-Agent` header in the format: `PLATFORM-system-code/version`
* `TransactionIDFromContext` will attempt to retrieve the transaction ID from the `*http.Request` context, and if it finds one will set the `X-Request-Id` header.

**Important** Neither extension will modify the request, it will only add a new header if none is already set.

# Usage

Create a new `*http.Client` which sets a `User-Agent` of `PLATFORM-system-code/version` for all requests:

```
trans := transport.NewTransport().WithStandardUserAgent("PLATFORM", "system-code")
client := &http.Client{Transport: trans}
```

`NewTransport()` assumes you would like to use the default `http.DefaultRoundTripper`, and the `TransactionIDFromContext` extension.

Automatically set the `X-Request-Id` (see the [Transaction ID Utils Go](https://github.com/Financial-Times/transactionid-utils-go) library for more details):

```
ctx := tidutils.TransactionAwareContext(context.TODO(), "tid_1234")
req, err := http.NewRequest("GET", uri, nil)
client.Do(req.WithContext(ctx))
```
