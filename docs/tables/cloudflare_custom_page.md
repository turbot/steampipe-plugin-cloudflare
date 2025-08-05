---
title: "Steampipe Table: cloudflare_custom_page - Query Cloudflare Custom Pages using SQL"
description: "Allows users to query Cloudflare Custom Pages, providing access to custom error or status page configurations, including page IDs, descriptions, states, URLs, creation and modification timestamps, preview targets, and required tokens at account or zone levels."
---

# Table: cloudflare_custom_page - Query Cloudflare Custom Pages using SQL

Custom Pages allow Cloudflare users to define custom error and status pages for display to end users. These pages support templated tokens and preview targets, and can be scoped at either the account or zone level.

## Table Usage Guide

The `cloudflare_custom_page` table provides insight into configurable custom page definitions within Cloudflare. As a security administrator or DevOps engineer, you can explore page-specific details, including page ID, description, state, URL, required tokens, preview target, creation timestamp and last modification timestamp. Use it to audit page customizations, monitor outdated error pages, and verify token usage for dynamic content rendering across account or zone contexts.

**Important Notes**
- You must specify either `account_id` or `zone_id` in a `where` or `join` clause to query this table.

## Examples

### Query all custom pages for a zone/account
```sql+postgres
select
  id,
  description,
  state,
  url,
  created_on,
  modified_on
from
  cloudflare_custom_page
where
  zone_id    = 'YOUR_ZONE_ID';
```

```sql+sqlite
select
  id,
  description,
  state,
  url,
  created_on,
  modified_on
from
  cloudflare_custom_page
where
  zone_id    = 'YOUR_ZONE_ID';
```

```sql+postgres
select
  id,
  description,
  state,
  url,
  created_on,
  modified_on
from
  cloudflare_custom_page
where
  account_id    = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  id,
  description,
  state,
  url,
  created_on,
  modified_on
from
  cloudflare_custom_page
where
  account_id    = 'YOUR_ACCOUNT_ID';
```

### Get a specific custom page by ID
```sql+postgres
select
  id,
  description,
  state,
  url,
  preview_target,
  required_tokens
from
  cloudflare_custom_page
where
  id = 'CUSTOM_PAGE_ID'
  account_id = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  id,
  description,
  state,
  url,
  preview_target,
  required_tokens
from
  cloudflare_custom_page
where
  id = 'CUSTOM_PAGE_ID'
  account_id = 'YOUR_ACCOUNT_ID';
```

### Query all custom pages recently created
```sql+postgres
select
  id,
  description,
  state,
  created_on
from
  cloudflare_custom_page
where
  account_id = 'YOUR_ACCOUNT_ID'
  and created_on >= now() - interval '7 days'
order by
  created_on desc;
```

```sql+sqlite
SELECT
  id,
  description,
  state,
  created_on
FROM
  cloudflare_custom_page
WHERE
  account_id = 'YOUR_ACCOUNT_ID'
  AND datetime(created_on) >= datetime('now', '-7 days')
ORDER BY
  created_on DESC;
```

### Query all customized error pages
```sql+postgres
select
  id,
  description,
  state,
  url
from
  cloudflare_custom_page
where
  account_id = 'YOUR_ACCOUNT_ID'
  and state == 'customized';
```

```sql+sqlite
select
  id,
  description,
  state,
  url
from
  cloudflare_custom_page
where
  account_id = 'YOUR_ACCOUNT_ID'
  and state == 'customized';
```
