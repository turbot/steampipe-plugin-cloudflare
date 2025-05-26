# Cloudflare Go API Structure Guide

## Overview

This guide explains how to explore and understand the Cloudflare Go API structure, including how to use Go documentation tools and understand API responses.

## Core API Structure

The Cloudflare Go API is organized into several main components:

### 1. Client Structure

```go
// Create a new client
client := cloudflare.NewClient(
    option.WithAPIKey("your-api-key"),
    option.WithAPIEmail("your-email"),
)
```

### 2. Main Resource Categories

The API is organized into logical groupings that match Cloudflare's product structure:

1. **Zones** - Domain and DNS management
2. **Accounts** - Account management and settings
3. **Access** - Zero Trust access control
4. **Load Balancers** - Traffic distribution
5. **Workers** - Serverless computing
6. **R2** - Object storage
7. **Firewall** - Security rules and settings

## Exploring the API

### 1. Using Go Doc

To explore the API structure using Go's documentation tools:

```bash
# View package documentation
go doc github.com/cloudflare/cloudflare-go/v4

# View specific type documentation
go doc github.com/cloudflare/cloudflare-go/v4.Zone

# View method documentation
go doc github.com/cloudflare/cloudflare-go/v4.Client.ListZones
```

### 2. Key Resource Types

#### Zones
```go
type Zone struct {
    ID                  string    // Zone identifier
    Name                string    // Domain name
    Status              string    // Status of the zone
    Type                string    // Zone type (full, partial)
    NameServers        []string  // Assigned nameservers
    OriginalNameServers []string  // Original nameservers
    Paused             bool      // If zone is paused
    CreatedOn          time.Time // Creation timestamp
    ModifiedOn         time.Time // Last modification timestamp
    // ... other fields
}
```

#### DNS Records
```go
type DNSRecord struct {
    ID        string    // Record identifier
    Type      string    // Record type (A, AAAA, CNAME, etc)
    Name      string    // DNS record name
    Content   string    // DNS record content
    ZoneID    string    // Zone identifier
    ZoneName  string    // Zone name
    CreatedOn time.Time // Creation timestamp
    ModifiedOn time.Time // Last modification timestamp
    Proxied   bool      // If record is proxied
    TTL       int       // Time to live
    // ... other fields
}
```

## Understanding API Responses

### 1. Response Structure

All API responses follow a common pattern:

```go
type Response struct {
    Result     interface{} // The actual data
    Success    bool        // If the request was successful
    Errors     []Error    // Any errors that occurred
    Messages   []Message  // Any messages from the API
    ResultInfo ResultInfo // Pagination information
}

type ResultInfo struct {
    Page       int // Current page
    PerPage    int // Items per page
    Count      int // Total items
    TotalCount int // Total available items
}
```

### 2. Error Handling

```go
type Error struct {
    Code    int    // Error code
    Message string // Error message
}
```

## Common API Operations

### 1. List Operations

Most list operations support pagination and filtering:

```go
// List zones with pagination
zones, err := client.ListZones(context.Background(), cloudflare.ListZoneParams{
    Page:     1,
    PerPage:  20,
    Name:     "example.com",
})

// Auto-paging for all results
iter := client.Zones.ListAutoPaging(context.TODO(), zones.ZoneListParams{})
for iter.Next() {
    zone := iter.Current()
    // Process zone
}
```

### 2. Create Operations

Create operations typically require a params struct:

```go
record, err := client.CreateDNSRecord(context.Background(), cloudflare.CreateDNSRecordParams{
    ZoneID: "zone-id",
    Type:   "A",
    Name:   "example.com",
    Content: "192.0.2.1",
    TTL:    3600,
})
```

### 3. Update Operations

Update operations usually require the resource ID and updated fields:

```go
err := client.UpdateDNSRecord(context.Background(), cloudflare.UpdateDNSRecordParams{
    ZoneID: "zone-id",
    ID:     "record-id",
    Type:   "A",
    Name:   "example.com",
    Content: "192.0.2.2",
})
```

## Best Practices

1. **Use Context**: Always pass a context to API calls for timeout and cancellation
2. **Handle Pagination**: Use ListAutoPaging for large result sets
3. **Error Handling**: Always check both the error return and Response.Success
4. **Rate Limiting**: Be aware of API rate limits and use appropriate retries
5. **Field Validation**: Use the Field type system for request parameters

## Available Resources

Here are the main resources available through the API:

1. **Zone Management**
   - List/Create/Update/Delete zones
   - Zone settings
   - DNS records
   - SSL/TLS configuration

2. **Account Management**
   - Account details
   - Account members
   - Account roles
   - Account settings

3. **Security**
   - Firewall rules
   - Access applications
   - Access policies
   - WAF rules

4. **Performance**
   - Load balancers
   - Page rules
   - Cache configuration
   - Workers

5. **Storage**
   - R2 buckets
   - R2 objects
   - Workers KV

## Exploring API Details

To get detailed information about specific API endpoints:

1. Use the [Cloudflare API documentation](https://api.cloudflare.com/)
2. Use Go's documentation tools
3. Look at the [example code](https://github.com/cloudflare/cloudflare-go/tree/master/examples)
4. Check the [Cloudflare Developer Docs](https://developers.cloudflare.com/api/) 