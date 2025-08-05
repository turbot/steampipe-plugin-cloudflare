---
title: "Steampipe Table: cloudflare_logpush_job - Query Cloudflare Logpush Jobs using SQL"
description: "Allows users to query Cloudflare Logpush Jobs, exposing log export settings including job ID, dataset, destination config, enabled status, last completion or error timestamps, error messages, output options, upload size/interval/record limits, and job name at account or zone level."
---

# Table: cloudflare_logpush_job - Query Cloudflare Logpush Jobs using SQL

Logpush Jobs automate near‑real‑time export of Cloudflare logs to an external destination, such as cloud storage. Supports datasets like HTTP requests or Workers Trace Events.

## Table Usage Guide

The `cloudflare_logpush_job` table provides insights into configured log shipping jobs within Cloudflare. As a security administrator or DevOps engineer, you can can review job ID, dataset, destination configuration, enablement flag, timestamps of last successful or failed runs, error messages, output options and thresholds (such as max upload bytes, interval or record counts), and job names. Use it to monitor job health, validate thresholds and destinations, track failed exports, and manage log delivery configuration across account or zone scope. 

**Important Notes**
- You must specify either `account_id` or `zone_id` in a `where` or `join` clause to query this table.

## Examples

### Query all logpush jobs for a zone/account
```sql+postgres
select
  id,
  name,
  dataset,
  destination_conf,
  enabled,
  last_complete,
  last_error
from
  cloudflare_logpush_job
where
  zone_id    = 'YOUR_ZONE_ID';
```

```sql+sqlite
select
  id,
  name,
  dataset,
  destination_conf,
  enabled,
  last_complete,
  last_error
from
  cloudflare_logpush_job
where
  zone_id    = 'YOUR_ZONE_ID';
```

```sql+postgres
select
  id,
  name,
  dataset,
  destination_conf,
  enabled,
  last_complete,
  last_error
from
  cloudflare_logpush_job
where
  account_id    = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  id,
  name,
  dataset,
  destination_conf,
  enabled,
  last_complete,
  last_error
from
  cloudflare_logpush_job
where
  account_id    = 'YOUR_ACCOUNT_ID';
```

### Get a specific logpush job
```sql+postgres
select
  id,
  name,
  dataset,
  enabled,
  destination_conf,
  output_options,
  max_upload_bytes,
  max_upload_records,
  max_upload_interval_seconds,
  error_message,
  last_complete,
  last_error
from
  cloudflare_logpush_job
where
  id = 123456789
  account_id = 'YOUR_ACCOUNT_ID' ;
```

```sql+sqlite
select
  id,
  name,
  dataset,
  enabled,
  destination_conf,
  output_options,
  max_upload_bytes,
  max_upload_records,
  max_upload_interval_seconds,
  error_message,
  last_complete,
  last_error
from
  cloudflare_logpush_job
where
  id = 123456789
  account_id = 'YOUR_ACCOUNT_ID' ;
```

### Query all logpush jobs with a recent failure
```sql+postgres
select
  id,
  name,
  enabled,
  last_error,
  error_message
from
  cloudflare_logpush_job
where
  account_id = 'YOUR_ACCOUNT_ID'
  and enabled = true
  and last_error is not null
order by
  last_error desc;
```

```sql+sqlite
select
  id,
  name,
  enabled,
  last_error,
  error_message
from
  cloudflare_logpush_job
where
  account_id = 'YOUR_ACCOUNT_ID'
  and enabled = true
  and last_error is not null
order by
  last_error desc;
```

### Query all disabled logpush jobs
```sql+postgres
select
  id,
  name,
  dataset,
  enabled
from
  cloudflare_logpush_job
where
  account_id = 'YOUR_ACCOUNT_ID'
  and enabled = false;
```

```sql+sqlite
select
  id,
  name,
  dataset,
  enabled
from
  cloudflare_logpush_job
where
  account_id = 'YOUR_ACCOUNT_ID'
  and enabled = false;
```

### Query all logpush jobs sending firewall events
```sql+postgres
select
  id,
  name,
  dataset
from
  cloudflare_logpush_job
where
  account_id = 'YOUR_ACCOUNT_ID'
  and dataset = 'firewall_events';
```

```sql+sqlite
select
  id,
  name,
  dataset
from
  cloudflare_logpush_job
where
  account_id = 'YOUR_ACCOUNT_ID'
  and dataset = 'firewall_events';
```
