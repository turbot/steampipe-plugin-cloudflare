---
title: "Steampipe Table: cloudflare_healthcheck - Query Cloudflare Healthchecks using SQL"
description: "Allows users to query Cloudflare Health Checks, surfacing configuration for origin monitoring including check ID, name, address, protocol, status, fail/success thresholds, intervals, retries, timeouts, suspension flag, and region settings at zone level."
---

# Table: cloudflare_healthcheck - Query Cloudflare Healthchecks using SQL

Health Checks allow Cloudflare to monitor origin server availability via scheduled tests from its edge network. Supports various protocols (HTTP, HTTPS, TCP) and region-based health verification.

## Table Usage Guide

The `cloudflare_healthcheck` table provides insights into health check definitions per zone within Cloudflare. As a security administrator or DevOps engineer, you can explore healthcheck ID, name, address, protocol type, thresholds for consecutive failures/successes, interval and timeout settings, retry counts, current health status, suspension flag, and the list of regions where checks are executed. Use it to audit active health monitoring, detect unhealthy origins, adjust thresholds, and verify regional test coverage.

**Important Notes**
- You must specify a `zone_id` in a `where` or `join` clause to query this table.

## Examples

### Query all healthcheck for a zone
Retrieves all healthchecks associated with a specific zone ID. Healthchecks are used to monitor the availability and performance of backend resources (e.g., servers or services).

```sql+postgres
select
  id,
  name,
  address,
  status,
  created_on,
  description
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID';
```

```sql+sqlite
select
  id,
  name,
  address,
  status,
  created_on,
  description
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID';
```

### Get a specific healthcheck by ID
Retrieves detailed information about a specific healthcheck, identified by its ID and the zone ID.

```sql+postgres
select
  id,
  name,
  address,
  status,
  failure_reason,
  consecutive_fails,
  consecutive_successes,
  interval,
  timeout,
  retries,
  suspended,
  modified_on
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID'
  and id = 'HEALTHCHECK_ID';
```

```sql+sqlite
select
  id,
  name,
  address,
  status,
  failure_reason,
  consecutive_fails,
  consecutive_successes,
  interval,
  timeout,
  retries,
  suspended,
  modified_on
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID'
  and id = 'HEALTHCHECK_ID';
```

### Query all unhealthy healthcheck with more than 'n' consectives fails
Retrieves all unhealthy healthchecks for a specific zone that have experienced 3 or more consecutive failures. It's useful for identifying problem areas and addressing persistent back-end service issues.

```sql+postgres
select
  id,
  name,
  status,
  failure_reason,
  consecutive_fails
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID'
  and consecutive_fails >= 3
  and status = 'unhealthy'
order by
  consecutive_fails desc;
```

```sql+sqlite
select
  id,
  name,
  status,
  failure_reason,
  consecutive_fails
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID'
  and consecutive_fails >= 3
  and status = 'unhealthy'
order by
  consecutive_fails desc;
```

### Query all suspended healthcheck
Retrieves all suspended healthchecks for a specific zone. A suspended healthcheck is one that has been temporarily paused and is not actively running probes.

```sql+postgres
select
  id,
  name,
  address,
  status,
  suspended
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID'
  and suspended = true;
```

```sql+sqlite
select
  id,
  name,
  address,
  status,
  suspended
from
  cloudflare_healthcheck
where
  zone_id = 'YOUR_ZONE_ID'
  and suspended = true;
```
