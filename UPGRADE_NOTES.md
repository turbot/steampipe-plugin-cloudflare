# Cloudflare Go SDK Upgrade Notes (v0.27.0 to v4.2.0)

## Overview
This document outlines the necessary changes required to upgrade from cloudflare-go v0.27.0 to v4.2.0. The upgrade includes several breaking changes and new features that need to be considered.

## Major Breaking Changes

### 1. Go Version Requirement
- Minimum Go version required is now 1.18+

### 2. Package Import Path
```go
// Old import
import "github.com/cloudflare/cloudflare-go"

// New import 
import "github.com/cloudflare/cloudflare-go/v4"
```

### 3. API Changes

#### Access Service
- The Access service has been moved to the `zero_trust` package
- Access functionality is now under `zero_trust.AccessService`
- Access applications are now managed through `zero_trust.AccessService.Applications`
- Access policies are now managed through `zero_trust.AccessService.Applications.Policies`
- The API methods have been reorganized to better reflect the resource hierarchy
- Example: `client.AccessPolicies()` is now `zero_trust.NewAccessService().Applications.Policies`

#### API Token Methods
- The API token methods have been reorganized in v4.2.0
- The old `client.APITokens()` and `client.UserAPITokens()` methods are no longer available
- Need to investigate the correct replacement methods in the new SDK structure

#### Page Rules Methods
- The page rules package structure has changed in v4.2.0
- The old types and methods are no longer available
- Need to investigate the correct package structure and methods

#### Load Balancer Methods
- The `load_balancing` package is no longer available
- Load balancer functionality has been moved to a new package
- Need to investigate the correct package structure and methods

#### Account Member Methods
- The `account` package structure has changed in v4.2.0
- Need to investigate the correct package structure and methods

### 4. Common Patterns
- Services are now created using `NewXService(opts...)`
- Methods take structured params with `cloudflare.F()` wrappers
- Options are passed using `option.RequestOption`
- Service methods follow a more consistent pattern across the SDK
- Many methods now require explicit parent resource IDs as parameters
- Example: `ListPolicies(ctx, appID, params)` instead of `params.ApplicationID`

### 5. Migration Strategy
1. Update import paths to use v4
2. Replace old service instantiation with new patterns
3. Update method calls to use new package structure
4. Add required parent resource IDs as explicit parameters
5. Update parameter structs to match new SDK requirements
6. Test each component after migration

#### Request Fields
- All request parameters are now wrapped in a generic `Field` type
- Use helpers like `String()`, `Int()`, `Float()`, or the generic `F[T]()` to construct fields
- Use `Null[T]()` to send null values
- Use `Raw[T](any)` for non-conforming values

Example:
```go
// Old way
params := FooParams{
    Name: "hello",
    Description: nil,
}

// New way
params := FooParams{
    Name: cloudflare.F("hello"),
    Description: cloudflare.Null[string](),
    Point: cloudflare.F(cloudflare.Point{
        X: cloudflare.Int(0),
        Y: cloudflare.Int(1),
    }),
}
```

#### Response Objects
- All fields in response structs are now value types (not pointers)
- Null or missing fields will have zero values
- New `.JSON` field available on responses for detailed property information
- Access to raw response data through `ExtraFields` map

Example:
```go
if res.Name == "" {
    // Check if name was null or missing
    isNull := res.JSON.Name.IsNull()
    isMissing := res.JSON.Name.IsMissing()
    
    // Access unexpected fields
    extraData := res.JSON.ExtraFields["unexpected_field"].Raw()
}
```

### 6. Error Handling
- API errors now return type `*cloudflare.Error`
- Use `errors.As` pattern for error handling

Example:
```go
_, err := client.Zones.Get(context.TODO(), zones.ZoneGetParams{
    ZoneID: cloudflare.F("023e105f4ecef8ad9ca31a8372d0c353"),
})
if err != nil {
    var apierr *cloudflare.Error
    if errors.As(err, &apierr) {
        println(string(apierr.DumpRequest(true)))
        println(string(apierr.DumpResponse(true)))
    }
    panic(err.Error())
}
```

## New Features

### 1. Request Options
- New functional options pattern using `option` package
- Options can be supplied at client or request level

Example:
```go
client := cloudflare.NewClient(
    option.WithHeader("X-Some-Header", "custom_header_info"),
)

client.Zones.New(context.TODO(), ...,
    option.WithHeader("X-Some-Header", "override_header"),
    option.WithJSONSet("some.json.path", map[string]string{"my": "object"}),
)
```

### 2. Pagination
- New `.ListAutoPaging()` methods for automatic pagination
- Simple `.List()` methods for single page fetching

Example:
```go
// Auto-paging
iter := client.Accounts.ListAutoPaging(context.TODO(), accounts.AccountListParams{})
for iter.Next() {
    account := iter.Current()
    fmt.Printf("%+v\n", account)
}

// Manual paging
page, err := client.Accounts.List(context.TODO(), accounts.AccountListParams{})
for page != nil {
    for _, account := range page.Result {
        fmt.Printf("%+v\n", account)
    }
    page, err = page.GetNextPage()
}
```

### 3. Timeouts
- Use context for request lifecycle timeout
- Use `option.WithRequestTimeout()` for per-retry timeout

Example:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
client.Zones.Edit(
    ctx,
    zones.ZoneEditParams{
        ZoneID: cloudflare.F("023e105f4ecef8ad9ca31a8372d0c353"),
    },
    option.WithRequestTimeout(20*time.Second),
)
```

### 4. File Uploads
- File upload parameters typed as `param.Field[io.Reader]`
- Use `cloudflare.FileParam()` helper for custom filename and content-type

Example:
```go
// File from filesystem
file, err := os.Open("/path/to/file")
params := api_gateway.UserSchemaNewParams{
    ZoneID: cloudflare.F("023e105f4ecef8ad9ca31a8372d0c353"),
    File: cloudflare.F[io.Reader](file),
}

// File from string with custom name/type
params := api_gateway.UserSchemaNewParams{
    ZoneID: cloudflare.F("023e105f4ecef8ad9ca31a8372d0c353"),
    File: cloudflare.FileParam(
        strings.NewReader(`{"hello": "foo"}`),
        "file.json",
        "application/json"
    ),
}
```

### 5. Retries
- Automatic retry of certain errors (2 retries by default)
- Configurable with `WithMaxRetries` option

Example:
```go
// Global config
client := cloudflare.NewClient(
    option.WithMaxRetries(0), // disable retries
)

// Per-request config
client.Zones.Get(
    context.TODO(),
    zones.ZoneGetParams{
        ZoneID: cloudflare.F("023e105f4ecef8ad9ca31a8372d0c353"),
    },
    option.WithMaxRetries(5),
)
```

## Migration Steps

1. Update Go version to 1.18+ if not already done
2. Update import paths to use v4
3. Convert all request parameter fields to use the new Field type system
4. Update error handling code to use the new Error type
5. Review and update any custom request/response handling
6. Test thoroughly, especially around error cases and parameter handling

## Additional Resources

- [Official Cloudflare Go SDK Documentation](https://github.com/cloudflare/cloudflare-go)
- [API Reference](https://developers.cloudflare.com/api/) 