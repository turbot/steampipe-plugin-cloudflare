---
title: "Steampipe Table: cloudflare_worker_script - Query Cloudflare Worker Scripts using SQL"
description: "Allows users to query Cloudflare Worker Scripts, providing metadata on deployed Workers including script ID, timestamps, usage model, presence of assets or modules, logpush flag, account association, subdomain availability, tail consumer list, and placement settings."
---

# Table: cloudflare_worker_script - Query Cloudflare Worker Scripts using SQL

Worker Scripts host custom serverless code executed at Cloudflare’s edge for enhanced logic, routing, and performance optimizations. They may include assets, modules, and optionally support logging via Logpush.

## Table Usage Guide

The `cloudflare_worker_script` table provides insights into metadata for Workers deployed per account within Cloudflare. As a security administrator or DevOps engineer, you can inspect script ID, creation and last modification timestamps, usage model, boolean flags for assets, modules and logpush enablement, JSON representing workers.dev subdomain mappings, list of log tail consumers, and placement configurations. Use it to audit Worker deployments, verify logging settings, monitor asset/module usage, check subdomain exposure, and review placement strategies across accounts.


## Examples

### Query all worker scripts
```sql+postgres
select
  id,
  account_id,
  account_name,
  created_on,
  modified_on,
  usage_model,
  has_assets,
  has_modules,
  logpush
from
  cloudflare_worker_script
```

```sql+sqlite
select
  id,
  account_id,
  account_name,
  created_on,
  modified_on,
  usage_model,
  has_assets,
  has_modules,
  logpush
from
  cloudflare_worker_script
```

### Query all worker scripts with worker.dev subdomain available
```sql+postgres
select
  id,
  account_id,
  account_name,
  subdomain
from
  cloudflare_worker_script
where
  subdomain ->> 'enabled' = 'false';
```

```sql+sqlite
select
  id,
  account_id,
  account_name,
  subdomain
from
  cloudflare_worker_script
where
  subdomain ->> 'enabled' = 'false';
```

### Query all worker scripts which have modules or assets
```sql+postgres
select
  id,
  account_name,
  has_assets,
  has_modules
from
  cloudflare_worker_script
where
  has_assets = true or has_modules = true;
```

```sql+sqlite
select
  id,
  account_name,
  has_assets,
  has_modules
from
  cloudflare_worker_script
where
  has_assets = true or has_modules = true;
```
