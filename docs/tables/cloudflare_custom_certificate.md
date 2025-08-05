---
title: "Steampipe Table: cloudflare_custom_certificate - Query Cloudflare Custom Certificates using SQL"
description: "Allows users to query Custom Certificates in Cloudflare, offering visibility into user-managed SSL/TLS certificates for Business and Enterprise plans, covering details such as hosts, issuer, expiration, upload and modification timestamps, bundle method, priority, status, and geo-key policies at the zone level."
---

# Table: cloudflare_custom_certificate - Query Cloudflare Custom Certificates using SQL

Custom Certificates allow Business and Enterprise customers to bring and manage their own SSL/TLS certificates within Cloudflare. These user‑provided certificates require manual lifecycle management—uploading, renewing, and monitoring expiration. Cloudflare organizes them into certificate packs (covering multiple signature algorithms per hostname group), and supports geo‑key restriction for private key locality.

## Table Usage Guide

The `cloudflare_custom_certificate` table provides insights into user-managed SSL/TLS certificates within Cloudflare. As a security administrator or DevOps engineer, you can explore certificate-specific details through this table, including hosts, issuer, expiration, upload and modification timestamps, bundle method, priority, status, and geo-key policy configurations—all scoped at the zone level. Leverage it to audit custom certificate deployments, monitor expiration lifecycles, manage overlapping certificate priorities, and enforce geo‑restriction policies across your Cloudflare infrastructure.

**Important Notes**
- You must specify a `zone_id` in a `where` or `join` clause to query this table.

## Examples

### Query all custom certificates for a zone
```sql+postgres
select
  id,
  hosts,
  issuer,
  status,
  expires_on
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID';
```

```sql+sqlite
select
  id,
  hosts,
  issuer,
  status,
  expires_on
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID';
```

### Get a specific custom certificate by its ID
```sql+postgres
select
  id,
  hosts,
  bundle_method,
  priority,
  uploaded_on,
  modified_on,
  geo_restrictions,
  policy
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID'
  and id      = 'CERTIFICATE_ID';
```

```sql+sqlite
select
  id,
  hosts,
  bundle_method,
  priority,
  uploaded_on,
  modified_on,
  geo_restrictions,
  policy
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID'
  and id      = 'CERTIFICATE_ID';
```

### Query all custom certificates expiring in the next 30 days for a zone
```sql+postgres
select
  id,
  hosts,
  issuer,
  expires_on
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID'
  and expires_on < now() + interval '30 days'
order by
  expires_on;
```

```sql+sqlite
select
  id,
  hosts,
  issuer,
  expires_on
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID'
  and expires_on < now() + interval '30 days'
order by
  expires_on;
```

### Order the custom certificates for a zone by priority
```sql+postgres
select
  id,
  hosts,
  issuer,
  priority,
  status
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID'
order by
  priority::integer asc;
```

```sql+sqlite
select
  id,
  hosts,
  issuer,
  priority,
  status
from
  cloudflare_custom_certificate
where
  zone_id = 'YOUR_ZONE_ID'
order by
  priority::integer asc;
```