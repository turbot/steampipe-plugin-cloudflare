---
title: "Steampipe Table: cloudflare_custom_certificate - Query Cloudflare Custom Certificates using SQL"
description: "Allows users to query Custom Certificates in Cloudflare, offering visibility into user-managed SSL/TLS certificates for Business and Enterprise plans, covering details such as hosts, issuer, expiration, upload and modification timestamps, bundle method, priority, status, and geo-key policies at the zone level."
---

# Table: cloudflare_custom_certificate - Query Cloudflare Custom Certificates using SQL

Custom Certificates allow Business and Enterprise customers to bring and manage their own SSL/TLS certificates within Cloudflare. These user‑provided certificates require manual lifecycle management—uploading, renewing, and monitoring expiration. Cloudflare organizes them into certificate packs (covering multiple signature algorithms per hostname group), and supports geo‑key restriction for private key locality.

## Table Usage Guide

The `cloudflare_custom_certificate` table provides insights into user-managed SSL/TLS certificates within Cloudflare. As a security administrator or DevOps engineer, you can explore certificate-specific details through this table, including hosts, issuer, expiration, upload and modification timestamps, bundle method, priority, status, and geo-key policy configurations—all scoped at the zone level. Leverage it to audit custom certificate deployments, monitor expiration lifecycles, manage overlapping certificate priorities, and enforce geo‑restriction policies across your Cloudflare infrastructure.

## Examples

### Query all custom certificates

### Query all custom certificates for a zone
Retrieves all custom SSL certificates associated with a specific zone ID. Custom certificates are domain-specific SSL/TLS certificates uploaded for use in Cloudflare.

```sql+postgres
select
  cc.id,
  cc.hosts,
  cc.issuer,
  cc.status,
  cc.expires_on,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID';
```

```sql+sqlite
select
  cc.id,
  cc.hosts,
  cc.issuer,
  cc.status,
  cc.expires_on,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID';
```

### Get a specific custom certificate by its ID
Retrieves detailed information about a specific custom certificate, identified by its ID and the zone ID.

```sql+postgres
select
  cc.id,
  cc.hosts,
  cc.bundle_method,
  cc.priority,
  cc.uploaded_on,
  cc.modified_on,
  cc.geo_restrictions,
  cc.policy,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID'
  and cc.id = 'CERTIFICATE_ID';
```

```sql+sqlite
select
  cc.id,
  cc.hosts,
  cc.bundle_method,
  cc.priority,
  cc.uploaded_on,
  cc.modified_on,
  cc.geo_restrictions,
  cc.policy,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID'
  and cc.id = 'CERTIFICATE_ID';
```

### Query all custom certificates expiring in the next 30 days for a zone
Retrieves all custom certificates for a specific zone that will expire in the next 30 days. This query is useful for monitoring expiring certificates to renew them before they expire.

```sql+postgres
select
  cc.id,
  cc.hosts,
  cc.issuer,
  cc.expires_on,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID'
  and cc.expires_on < now() + interval '30 days'
order by
  cc.expires_on;
```

```sql+sqlite
select
  cc.id,
  cc.hosts,
  cc.issuer,
  cc.expires_on,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID'
  and cc.expires_on < datetime('now', '+30 days')
order by
  cc.expires_on;
```

### Order the custom certificates for a zone by priority
Retrieves all custom certificates for a specific zone and orders them by their priority. Numeric priority values determine which certificate is selected when multiple certificates apply to a hostname.

```sql+postgres
select
  cc.id,
  cc.hosts,
  cc.issuer,
  cc.priority,
  cc.status,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID'
order by
  cc.priority::integer asc;
```

```sql+sqlite
select
  cc.id,
  cc.hosts,
  cc.issuer,
  cc.priority,
  cc.status,
  z.name as zone_name
from
  cloudflare_custom_certificate cc
join
  cloudflare_zone z
on
  cc.zone_id = z.id
where
  cc.zone_id = 'YOUR_ZONE_ID'
order by
  cast(cc.priority as integer) asc;
```