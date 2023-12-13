---
title: "Steampipe Table: cloudflare_firewall_rule - Query Cloudflare Firewall Rules using SQL"
description: "Allows users to query Cloudflare Firewall Rules, specifically rules that control security and access to a site, providing insights into configurations and potential security risks."
---

# Table: cloudflare_firewall_rule - Query Cloudflare Firewall Rules using SQL

Cloudflare Firewall Rules are a security feature that determines which traffic you want to allow to your website. They are customizable and can be used to mitigate against threats, control access, and block traffic from certain IP addresses or regions. Firewall rules can be set up to match against incoming HTTP traffic, and actions can be taken based on the rule match.

## Table Usage Guide

The `cloudflare_firewall_rule` table provides insights into Firewall Rules within Cloudflare. As a security engineer, explore rule-specific details through this table, including rule configurations, action taken, and associated metadata. Utilize it to uncover information about rules, such as those blocking specific IP addresses or regions, and the verification of rule configurations.

## Examples

### Basic info
Explore the creation dates of specific firewall rules to gain insights into their historical context and assess potential patterns or anomalies. This may assist in troubleshooting or optimizing your firewall configuration.

```sql+postgres
select
  id,
  zone_id,
  created_on
from
  cloudflare_firewall_rule;
```

```sql+sqlite
select
  id,
  zone_id,
  created_on
from
  cloudflare_firewall_rule;
```

### List paused firewall rules
Discover the segments that have paused firewall rules. This can be useful for identifying potential security vulnerabilities or areas where firewall protection is currently inactive.

```sql+postgres
select
  id,
  zone_id,
  paused
from
  cloudflare_firewall_rule
where
  paused;
```

```sql+sqlite
select
  id,
  zone_id,
  paused
from
  cloudflare_firewall_rule
where
  paused = 1;
```

### List firewall rules that block requests based on IP reputation
Analyze firewall rules to understand which ones are set to block based on IP reputation, helping to enhance security by identifying potential threats. This is particularly useful in preventing access from high-risk IP addresses.

```sql+postgres
select
  id,
  zone_id,
  filter,
  action
from
  cloudflare_firewall_rule
where
  action = 'block'
  and filter ->> 'expression' = '(cf.threat_score gt 1)';
```

```sql+sqlite
select
  id,
  zone_id,
  filter,
  action
from
  cloudflare_firewall_rule
where
  action = 'block'
  and json_extract(filter, '$.expression') = '(cf.threat_score gt 1)';
```