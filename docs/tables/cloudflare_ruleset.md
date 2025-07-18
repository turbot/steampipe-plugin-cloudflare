---
title: "Steampipe Table: cloudflare_ruleset - Query Cloudflare Rulesets using SQL"
description: "Allows users to query Rulesets in Cloudflare, providing insights into security rules, WAF configurations, and custom rule configurations at both account and zone levels."
---

# Table: cloudflare_ruleset - Query Cloudflare Rulesets using SQL

Cloudflare Rulesets provide a powerful framework for configuring rules to process HTTP requests. Rulesets contain one or more rules that can be used to filter, block, allow, log, or otherwise transform HTTP requests. They are the foundation of Cloudflare's security and performance features including WAF, Page Rules, Rate Limiting, and more.

## Table Usage Guide

The `cloudflare_ruleset` table provides insights into rulesets within Cloudflare. As a security administrator or DevOps engineer, you can explore ruleset-specific details through this table, including rule configurations, phases, kinds, and their associated accounts or zones. Utilize it to uncover information about security policies, understand rule hierarchies, audit configurations, and manage firewall rules across your Cloudflare infrastructure.

**Important Notes**
- You must specify either `account_id` or `zone_id` in a `where` or `join` clause to query this table.

## Examples

### Query all rulesets for an account
Explore all rulesets associated with a specific account to understand the security policies and rule configurations in place. This is useful for security audits and compliance reviews.

```sql+postgres
select
  id,
  name,
  kind,
  phase,
  description,
  version
from
  cloudflare_ruleset
where
  account_id = 'your_account_id';
```

```sql+sqlite
select
  id,
  name,
  kind,
  phase,
  description,
  version
from
  cloudflare_ruleset
where
  account_id = 'your_account_id';
```

### Query all rulesets for a specific zone
Retrieve all rulesets configured for a particular zone to understand zone-specific security configurations and custom rules.

```sql+postgres
select
  id,
  name,
  kind,
  phase,
  description,
  last_updated,
  rules
from
  cloudflare_ruleset
where
  zone_id = 'your_zone_id';
```

```sql+sqlite
select
  id,
  name,
  kind,
  phase,
  description,
  last_updated,
  rules
from
  cloudflare_ruleset
where
  zone_id = 'your_zone_id';
```

### Get a specific ruleset by ID
Retrieve detailed information about a particular ruleset, including all its rules and configuration details.

```sql+postgres
select
  *
from
  cloudflare_ruleset
where
  id = 'ruleset_id'
  and account_id = 'your_account_id';
```

```sql+sqlite
select
  *
from
  cloudflare_ruleset
where
  id = 'ruleset_id'
  and account_id = 'your_account_id';
```

### List custom rulesets with rule counts
Identify custom rulesets and count the number of rules in each to understand rule complexity and management requirements.

```sql+postgres
select
  id,
  name,
  kind,
  phase,
  json_array_length(rules) as rule_count
from
  cloudflare_ruleset
where
  account_id = 'your_account_id'
  and kind = 'custom'
order by
  rule_count desc;
```

```sql+sqlite
select
  id,
  name,
  kind,
  phase,
  json_array_length(rules) as rule_count
from
  cloudflare_ruleset
where
  account_id = 'your_account_id'
  and kind = 'custom'
order by
  rule_count desc;
```

### List rulesets by phase
Explore rulesets organized by their execution phase to understand the order and timing of rule processing in your Cloudflare configuration.

```sql+postgres
select
  phase,
  count(*) as ruleset_count,
  array_agg(name) as ruleset_names
from
  cloudflare_ruleset
where
  account_id = 'your_account_id'
group by
  phase
order by
  phase;
```

```sql+sqlite
select
  phase,
  count(*) as ruleset_count,
  group_concat(name) as ruleset_names
from
  cloudflare_ruleset
where
  account_id = 'your_account_id'
group by
  phase
order by
  phase;
```

### Find recently updated rulesets
Identify rulesets that have been modified recently to track configuration changes and updates.

```sql+postgres
select
  id,
  name,
  kind,
  phase,
  last_updated,
  description
from
  cloudflare_ruleset
where
  account_id = 'your_account_id'
  and last_updated > now() - interval '7 days'
order by
  last_updated desc;
```

```sql+sqlite
select
  id,
  name,
  kind,
  phase,
  last_updated,
  description
from
  cloudflare_ruleset
where
  account_id = 'your_account_id'
  and last_updated > datetime('now', '-7 days')
order by
  last_updated desc;
```

### List all rulesets across all zones
Get a comprehensive view of all rulesets by joining zone-level rulesets to understand the complete security policy landscape.

```sql+postgres
select
  r.id,
  r.name,
  r.kind,
  r.phase,
  r.account_id,
  r.zone_id,
  z.name as zone_name
from
  cloudflare_zone z
left join
  cloudflare_ruleset r on r.zone_id = z.id;
```

```sql+sqlite
select
  r.id,
  r.name,
  r.kind,
  r.phase,
  r.account_id,
  r.zone_id,
  z.name as zone_name
from
  cloudflare_zone z
left join
  cloudflare_ruleset r on r.zone_id = z.id;
```

### List account-level rulesets with account information
Retrieve all account-level rulesets along with account details to understand organization-wide security policies and account configurations.

```sql+postgres
select
  r.id,
  r.name,
  r.kind,
  r.phase,
  r.description,
  r.version,
  a.name as account_name,
  a.type as account_type,
  a.settings
from
  cloudflare_account a
join
  cloudflare_ruleset r on r.account_id = a.id
order by
  a.name, r.phase, r.name;
```

```sql+sqlite
select
  r.id,
  r.name,
  r.kind,
  r.phase,
  r.description,
  r.version,
  a.name as account_name,
  a.type as account_type,
  a.settings
from
  cloudflare_account a
join
  cloudflare_ruleset r on r.account_id = a.id
order by
  a.name, r.phase, r.name;
```
