---
title: "Steampipe Table: cloudflare_managed_transform - Query Cloudflare Managed Transforms using SQL"
description: "Allows users to query Cloudflare Managed Transforms, surfacing configuration for HTTP request and response header modifications including transform ID, type, enabled status, conflict detection, and associated zone information."
---

# Table: cloudflare_managed_transform - Query Cloudflare Managed Transforms using SQL

Managed Transforms allow you to perform common adjustments to HTTP request and response headers with the click of a button. These pre-configured transforms help standardize header handling across your zones without requiring custom rules.

## Table Usage Guide

The `cloudflare_managed_transform` table provides insights into managed transform configurations per zone within Cloudflare. As a security administrator or DevOps engineer, you can explore transform ID, type (request or response header), enabled status, conflict detection flags, and the list of conflicting transforms. Use it to audit active header modifications, detect configuration conflicts, and verify transform settings across zones.

## Examples

### Query all managed transforms for a zone
Retrieves all managed transforms associated with a specific zone ID. Managed transforms are used to modify HTTP request and response headers automatically.

```sql+postgres
select
  mt.id,
  mt.type,
  mt.enabled,
  mt.has_conflict,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.zone_id = 'YOUR_ZONE_ID';
```

```sql+sqlite
select
  mt.id,
  mt.type,
  mt.enabled,
  mt.has_conflict,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.zone_id = 'YOUR_ZONE_ID';
```

### Query all enabled managed transforms
Retrieves all enabled managed transforms across all zones. This is useful for understanding which header modifications are currently active.

```sql+postgres
select
  mt.id,
  mt.type,
  mt.enabled,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.enabled = true;
```

```sql+sqlite
select
  mt.id,
  mt.type,
  mt.enabled,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.enabled = true;
```

### Query all managed transforms with conflicts
Retrieves all managed transforms that have conflicts with other transforms. This helps identify configuration issues that may need resolution.

```sql+postgres
select
  mt.id,
  mt.type,
  mt.enabled,
  mt.has_conflict,
  mt.conflicts_with,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.has_conflict = true;
```

```sql+sqlite
select
  mt.id,
  mt.type,
  mt.enabled,
  mt.has_conflict,
  mt.conflicts_with,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.has_conflict = true;
```

### Query all request header transforms
Retrieves all managed transforms that modify request headers. This is useful for auditing inbound traffic modifications.

```sql+postgres
select
  mt.id,
  mt.enabled,
  mt.has_conflict,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.type = 'request_header';
```

```sql+sqlite
select
  mt.id,
  mt.enabled,
  mt.has_conflict,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.type = 'request_header';
```

### Query all response header transforms
Retrieves all managed transforms that modify response headers. This is useful for auditing outbound traffic modifications.

```sql+postgres
select
  mt.id,
  mt.enabled,
  mt.has_conflict,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.type = 'response_header';
```

```sql+sqlite
select
  mt.id,
  mt.enabled,
  mt.has_conflict,
  z.name as zone_name
from
  cloudflare_managed_transform mt
join
  cloudflare_zone z
on
  mt.zone_id = z.id
where
  mt.type = 'response_header';
```

### Count managed transforms by type and status
Provides a summary of managed transforms grouped by type and enabled status. This gives an overview of your transform configuration across all zones.

```sql+postgres
select
  type,
  enabled,
  count(*) as transform_count
from
  cloudflare_managed_transform
group by
  type,
  enabled
order by
  type,
  enabled;
```

```sql+sqlite
select
  type,
  enabled,
  count(*) as transform_count
from
  cloudflare_managed_transform
group by
  type,
  enabled
order by
  type,
  enabled;
```

