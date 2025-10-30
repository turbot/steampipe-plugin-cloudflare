---
title: "Steampipe Table: cloudflare_dns_record - Query Cloudflare DNS Records using SQL"
description: "Allows users to query DNS Records in Cloudflare, providing insights into the configuration details, zone information, and other metadata associated with each DNS record."
---

# Table: cloudflare_dns_record - Query Cloudflare DNS Records using SQL

Cloudflare DNS is a service that provides global, fast, and secure Domain Name System services. It is built on a network that is unified, distributed, and operates in real-time, providing users with reliable and highly available DNS services. It also offers DNSSEC to protect against forged DNS answers and other security threats.

## Table Usage Guide

The `cloudflare_dns_record` table provides insights into DNS records within Cloudflare. As a network administrator, you can explore record-specific details through this table, including the type of record, associated zone, and configuration settings. Utilize it to uncover information about DNS records, such as those with certain configurations, the zones they are associated with, and their current status.

## Examples

### List all records from each zone
Explore which DNS records belong to each zone in your Cloudflare account. This allows you to understand the distribution and organization of your DNS records, aiding in efficient management and troubleshooting.

```sql+postgres
select
  *
from
  cloudflare_dns_record
```

```sql+sqlite
select
  *
from
  cloudflare_dns_record
```

### List DNS records filtered by type
Filter DNS records by their record type (e.g., A, CNAME, MX) for detailed management and troubleshooting of specific record types.

```sql+postgres
select
  name,
  type,
  content,
  ttl
from
  cloudflare_dns_record
where
  type = 'A';
```

```sql+sqlite
select
  name,
  type,
  content,
  ttl
from
  cloudflare_dns_record
where
  type = 'A';
```

### List DNS records that are proxied by Cloudflare
Explore DNS records that are utilizing Cloudflare's origin protection, helping optimize security and performance for those records.

```sql+postgres
select
  name,
  type,
  proxied
from
  cloudflare_dns_record
where
  proxied = true;
```

```sql+sqlite
select
  name,
  type,
  proxied
from
  cloudflare_dns_record
where
  proxied = 1;
```

### List DNS records for a zone with specific TTL values
Filter DNS records for a zone based on their ttl value, which is crucial for optimizing DNS cache times and performance.

```sql+postgres
select
  name,
  type,
  ttl
from
  cloudflare_dns_record
where
  zone_id = 'YOUR_ZONE_ID'
  and ttl > 300;
```

```sql+sqlite
select
  name,
  type,
  ttl
from
  cloudflare_dns_record
where
  zone_id = 'YOUR_ZONE_ID'
  and ttl > 300;
```

### List MX records in priority order
Determine the areas in which MX records are prioritized within a specific zone in your Cloudflare DNS, providing a clear order of priority. This is useful for managing mail exchange servers and ensuring smooth email delivery.

```sql+postgres
select
  name,
  type,
  priority
from
  cloudflare_dns_record
where
  zone_id = 'YOUR_ZONE_ID'
  and type = 'MX'
order by
  priority;
```

```sql+sqlite
select
  name,
  type,
  priority
from
  cloudflare_dns_record
where
  zone_id = 'YOUR_ZONE_ID'
  and type = 'MX'
order by
  priority;
```
