---
title: "Steampipe Table: cloudflare_api_token - Query Cloudflare API Tokens using SQL"
description: "Allows users to query Cloudflare API Tokens, providing insights into the security and access control of Cloudflare accounts."
---

# Table: cloudflare_api_token - Query Cloudflare API Tokens using SQL

Cloudflare API Tokens are a resource within Cloudflare that allows you to manage and control access to your Cloudflare accounts. They provide granular permissions and security for your Cloudflare operations. API Tokens are an essential part of managing access to Cloudflare resources and ensuring the security of your account.

## Table Usage Guide

The `cloudflare_api_token` table provides insights into the API Tokens within Cloudflare. As a DevOps engineer, explore token-specific details through this table, including permissions, associated policies, and metadata. Utilize it to uncover information about tokens, such as those with extensive permissions, the specific resources they can access, and the verification of associated policies.

**Important Notes**
- This table is only available when using `email` and `api_key` for authentication credentials. It does not work when using `api_token` for access.

## Examples

### Query all API tokens for this user
Explore all API tokens linked to a user's account to assess their validity and security. This can be beneficial in identifying any potential risks or outdated tokens that need to be updated or removed.

```sql+postgres
select
  *
from
  cloudflare_api_token;
```

```sql+sqlite
select
  *
from
  cloudflare_api_token;
```

### List API tokens by age, oldest first
Analyze the age of your API tokens to understand which ones have been in use the longest. This can help in managing token lifecycle and ensuring that older, potentially exposed tokens are refreshed or retired.

```sql+postgres
select
  id,
  name,
  status,
  issued_on,
  date_part('day', now() - issued_on) as age
from
  cloudflare_api_token
order by
  issued_on desc;
```

```sql+sqlite
select
  id,
  name,
  status,
  issued_on,
  julianday('now') - julianday(issued_on) as age
from
  cloudflare_api_token
order by
  issued_on desc;
```

### List API tokens expiring with the next 14 days
Explore active API tokens that are due to expire within the next two weeks. This can be useful for maintaining security and ensuring continuous service by renewing or replacing the tokens before they expire.

```sql+postgres
select
  id,
  name,
  status,
  expires_on
from
  cloudflare_api_token
where
  status is active
  and expires_on < current_date + interval '14 days';
```

```sql+sqlite
select
  id,
  name,
  status,
  expires_on
from
  cloudflare_api_token
where
  status = 'active'
  and expires_on < date('now', '+14 days');
```

### List all permissions granted to each API token
Explore the permissions associated with each API token to gain insights into their access levels and control mechanisms. This can be beneficial in maintaining security standards and ensuring appropriate access rights are granted.

```sql+postgres
select
  name,
  policy ->> 'effect',
  policy -> 'resources',
  perm_group ->> 'name'
from
  cloudflare_api_token,
  jsonb_array_elements(policies) as policy,
  jsonb_array_elements(policy -> 'permission_groups') as perm_group;
```

```sql+sqlite
select
  name,
  json_extract(policy.value, '$.effect'),
  policy.value,
  json_extract(perm_group.value, '$.name')
from
  cloudflare_api_token,
  json_each(policies) as policy,
  json_each(json_extract(policy.value, '$.permission_groups')) as perm_group;
```