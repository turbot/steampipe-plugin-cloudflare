---
title: "Steampipe Table: cloudflare_notification_policy - Query Cloudflare Notification Policies using SQL"
description: "Allows users to query Cloudflare Notification Policies, revealing alert configuration metadata such as policy ID, name, alert type, interval, filters, delivery mechanisms, enabled state, creation and modification timestamps at the account level."
---

# Table: cloudflare_notification_policy - Query Cloudflare Notification Policies using SQL

Notification Policies enable configurable alerts based on Cloudflare account-level events. They define what triggers alerts, frequency of re-alerting, filters for event inclusion, and channels for delivery.

## Table Usage Guide

The `cloudflare_notification_policy` table provides insights into notification policy definitions per account within Cloudflare. As a security administrator or DevOps engineer, you can review policy ID, name, alert type, re-alert interval, enabled state, description, JSON-encoded filters and mechanisms, and creation/modification timestamps. Use it to audit alert policies, identify disabled or misconfigured alerts, filter by event types, assess delivery mechanisms, and maintain event-driven monitoring.

**Important Notes**
- You must specify an `account_id` in a `where` or `join` clause to query this table.

## Examples

### Query all notification policies for an account
Retrieves all notification policies for a specific account ID. Notification policies define the types of alerts and notifications configured for an account.

```sql+postgres
select
  id,
  name,
  alert_type,
  alert_interval,
  enabled,
  created,
  modified
from
  cloudflare_notification_policy
where
  account_id = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  id,
  name,
  alert_type,
  alert_interval,
  enabled,
  created,
  modified
from
  cloudflare_notification_policy
where
  account_id = 'YOUR_ACCOUNT_ID';
```

### Get a specific notification policy by ID
Retrieves detailed information about a specific notification policy, identified by its ID and the account ID.

```sql+postgres
select
  id,
  name,
  description,
  alert_type,
  alert_interval,
  enabled,
  filters,
  mechanisms
from
  cloudflare_notification_policy
where
  id = 'NOTIFICATION_POLICY_ID'
  and account_id = 'YOUR_ACCOUNT_ID';
```

```sql+sqlite
select
  id,
  name,
  description,
  alert_type,
  alert_interval,
  enabled,
  filters,
  mechanisms
from
  cloudflare_notification_policy
where
  id = 'NOTIFICATION_POLICY_ID'
  and account_id = 'YOUR_ACCOUNT_ID';
```

### Query all disabled notification policies
Retrieves all disabled notification policies (enabled = false) for a specific account ID. Disabled policies represent configurations that are not actively generating alerts.

```sql+postgres
select
  id,
  name,
  alert_type,
  enabled
from
  cloudflare_notification_policy
where
  account_id = 'YOUR_ACCOUNT_ID'
  and enabled = false;
```

```sql+sqlite
select
  id,
  name,
  alert_type,
  enabled
from
  cloudflare_notification_policy
where
  account_id = 'YOUR_ACCOUNT_ID'
  and enabled = false;
```

### Query all notification policies for advanced DDoS alert
Retrieves all notification policies of type advanced_ddos_attack_l7_alert (specific to Layer 7 DDoS attack notifications) for a given account ID.

```sql+postgres
select
  id,
  name,
  alert_type,
  description
from
  cloudflare_notification_policy
where
  account_id = 'YOUR_ACCOUNT_ID'
  and alert_type = 'advanced_ddos_attack_l7_alert';
```

```sql+sqlite
select
  id,
  name,
  alert_type,
  description
from
  cloudflare_notification_policy
where
  account_id = 'YOUR_ACCOUNT_ID'
  and alert_type = 'advanced_ddos_attack_l7_alert';
```
