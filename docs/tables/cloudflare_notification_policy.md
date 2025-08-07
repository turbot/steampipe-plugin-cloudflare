---
title: "Steampipe Table: cloudflare_notification_policy - Query Cloudflare Notification Policies using SQL"
description: "Allows users to query Cloudflare Notification Policies, revealing alert configuration metadata such as policy ID, name, alert type, interval, filters, delivery mechanisms, enabled state, creation and modification timestamps at the account level."
---

# Table: cloudflare_notification_policy - Query Cloudflare Notification Policies using SQL

Notification Policies enable configurable alerts based on Cloudflare account-level events. They define what triggers alerts, frequency of re-alerting, filters for event inclusion, and channels for delivery.

## Table Usage Guide

The `cloudflare_notification_policy` table provides insights into notification policy definitions per account within Cloudflare. As a security administrator or DevOps engineer, you can review policy ID, name, alert type, re-alert interval, enabled state, description, JSON-encoded filters and mechanisms, and creation/modification timestamps. Use it to audit alert policies, identify disabled or misconfigured alerts, filter by event types, assess delivery mechanisms, and maintain event-driven monitoring.

## Examples

### Query all notification policies for an account
Retrieves all notification policies for a specific account ID. Notification policies define the types of alerts and notifications configured for an account.

```sql+postgres
select
  n.id,
  n.name,
  n.alert_type,
  n.alert_interval,
  n.enabled,
  n.created,
  n.modified,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.account_id = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  n.id,
  n.name,
  n.alert_type,
  n.alert_interval,
  n.enabled,
  n.created,
  n.modified,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.account_id = 'YOUR_ACCOUNT_ID';
```

### Get a specific notification policy by ID
Retrieves detailed information about a specific notification policy, identified by its ID and the account ID.

```sql+postgres
select
  n.id,
  n.name,
  n.description,
  n.alert_type,
  n.alert_interval,
  n.enabled,
  n.filters,
  n.mechanisms,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.id = 'NOTIFICATION_POLICY_ID'
  and n.account_id = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  n.id,
  n.name,
  n.description,
  n.alert_type,
  n.alert_interval,
  n.enabled,
  n.filters,
  n.mechanisms,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.id = 'NOTIFICATION_POLICY_ID'
  and n.account_id = 'YOUR_ACCOUNT_ID';
```

### Query all disabled notification policies
Retrieves all disabled notification policies (enabled = false) for a specific account ID. Disabled policies represent configurations that are not actively generating alerts.

```sql+postgres
select
  n.id,
  n.name,
  n.alert_type,
  n.enabled,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.account_id = 'YOUR_ACCOUNT_ID'
  and n.enabled = false;
```

```sql+sqlite
select
  n.id,
  n.name,
  n.alert_type,
  n.enabled,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.account_id = 'YOUR_ACCOUNT_ID'
  and n.enabled = false;
```

### Query all notification policies for advanced DDoS alert
Retrieves all notification policies of type advanced_ddos_attack_l7_alert (specific to Layer 7 DDoS attack notifications) for a given account ID.

```sql+postgres
select
  n.id,
  n.name,
  n.alert_type,
  n.description,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.account_id = 'YOUR_ACCOUNT_ID'
  and n.alert_type = 'advanced_ddos_attack_l7_alert';
```

```sql+sqlite
select
  n.id,
  n.name,
  n.alert_type,
  n.description,
  ca.name as account_name
from
  cloudflare_notification_policy n
join
  cloudflare_account ca
on
  n.account_id = ca.id
where
  n.account_id = 'YOUR_ACCOUNT_ID'
  and n.alert_type = 'advanced_ddos_attack_l7_alert';
```
