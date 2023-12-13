---
title: "Steampipe Table: cloudflare_page_rule - Query Cloudflare Page Rules using SQL"
description: "Allows users to query Cloudflare Page Rules, specifically the rules set up for a specific domain, providing insights into website traffic flow and potential anomalies."
---

# Table: cloudflare_page_rule - Query Cloudflare Page Rules using SQL

Cloudflare Page Rules is a feature within Cloudflare that allows you to customize your website's functionality and performance on a page-by-page basis. It provides a centralized way to set up and manage rules for various web pages, including redirection, caching, security, and more. Cloudflare Page Rules helps you optimize the delivery of your web content and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `cloudflare_page_rule` table provides insights into the rules set up for specific web pages within Cloudflare. As a DevOps engineer, explore rule-specific details through this table, including targets, actions, and associated metadata. Utilize it to uncover information about rules, such as those related to caching, security, or redirection, and the verification of their correct implementation.

## Examples

### Basic info
Explore which rules are active within your Cloudflare settings to prioritize their importance and manage your web traffic effectively. This can help optimize your website's performance and security.

```sql+postgres
select
  id,
  zone_id,
  status,
  priority
from
  cloudflare_page_rule;
```

```sql+sqlite
select
  id,
  zone_id,
  status,
  priority
from
  cloudflare_page_rule;
```

### List disabled page rules
Explore which page rules are currently disabled in your Cloudflare settings. This information can help you understand potential vulnerabilities or areas of your website that are not currently protected by active rules.

```sql+postgres
select
  id,
  zone_id,
  status
from
  cloudflare_page_rule
where
  status = 'disabled';
```

```sql+sqlite
select
  id,
  zone_id,
  status
from
  cloudflare_page_rule
where
  status = 'disabled';
```

### List page rules that do not have the Always Online feature enabled
Assess the elements within your website's page rules to identify those that lack the Always Online feature. This can be useful to ensure your site remains accessible even when your server goes offline, enhancing user experience and site reliability.

```sql+postgres
select
  id,
  zone_id,
  action ->> 'value' as always_online
from
  cloudflare_page_rule,
  jsonb_array_elements(actions) as action
where
  action ->> 'id' = 'always_online'
  and action ->> 'value' = 'off';
```

```sql+sqlite
select
  id,
  zone_id,
  json_extract(action.value, '$.value') as always_online
from
  cloudflare_page_rule,
  json_each(actions) as action
where
  json_extract(action.value, '$.id') = 'always_online'
  and json_extract(action.value, '$.value') = 'off';
```