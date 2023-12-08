---
title: "Steampipe Table: cloudflare_user_audit_log - Query Cloudflare User Audit Logs using SQL"
description: "Allows users to query User Audit Logs in Cloudflare, providing insights into user activities, changes made, and potential security issues."
---

# Table: cloudflare_user_audit_log - Query Cloudflare User Audit Logs using SQL

Cloudflare User Audit Log is a feature within Cloudflare that allows you to monitor and track user activities and changes within your Cloudflare account. It provides a comprehensive and searchable log of all actions performed by users, helping to identify potential security issues and ensure compliance. Cloudflare User Audit Log is an essential tool for maintaining visibility and control over your Cloudflare account.

## Table Usage Guide

The `cloudflare_user_audit_log` table provides insights into user activities and changes within Cloudflare. As a security analyst, explore detailed logs through this table, including the actions performed, the user who performed them, and the time of the action. Utilize it to monitor user behavior, identify potential security issues, and ensure compliance with security policies.

## Examples

### Basic info
Explore the recent changes made by users in your Cloudflare account. This query helps you monitor user activity, providing insights into what modifications were made, when, and by whom, enhancing your account's security and accountability.

```sql+postgres
select
  actor_email,
  actor_type,
  l.when,
  action ->> 'type' as action_type,
  action ->> 'result' as action_result,
  jsonb_pretty(new_value_json) as new_value,
  jsonb_pretty(old_value_json) as old_value,
  owner_id
from
  cloudflare_user_audit_log l;
```

```sql+sqlite
select
  actor_email,
  actor_type,
  l.when,
  json_extract(action, '$.type') as action_type,
  json_extract(action, '$.result') as action_result,
  new_value_json as new_value,
  old_value_json as old_value,
  owner_id
from
  cloudflare_user_audit_log l;
```

### Get all the users' activities in the last 10 days
Explore the recent activities of all users in the past 10 days. This allows for a comprehensive review of user actions, helping to identify any unusual patterns or potential security concerns.

```sql+postgres
select
  actor_email,
  actor_type,
  l.when,
  action ->> 'type' as action_type,
  action ->> 'result' as action_result,
  jsonb_pretty(new_value_json) as new_value,
  jsonb_pretty(old_value_json) as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  l.when > now() - interval '10' day;
```

```sql+sqlite
select
  actor_email,
  actor_type,
  l.when,
  json_extract(action, '$.type') as action_type,
  json_extract(action, '$.result') as action_result,
  new_value_json as new_value,
  old_value_json as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  l.when > datetime('now','-10 day');
```

### Get all the users' activities for a particular timeline
Explore the user activity within a specific timeframe to gain insights into their actions and results. This is beneficial for auditing purposes, understanding user behavior, and identifying any unusual or suspicious activity.

```sql+postgres
select
  actor_email,
  actor_type,
  l.when,
  action ->> 'type' as action_type,
  action ->> 'result' as action_result,
  jsonb_pretty(new_value_json) as new_value,
  jsonb_pretty(old_value_json) as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  l.when > '2023-06-04' and l.when < '2023-06-07';
```

```sql+sqlite
select
  actor_email,
  actor_type,
  l.when,
  json_extract(action, '$.type') as action_type,
  json_extract(action, '$.result') as action_result,
  new_value_json as new_value,
  old_value_json as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  l.when > '2023-06-04' and l.when < '2023-06-07';
```

### Get all the activities of a particular user
Explore the actions of a specific user to gain insights into their activities and changes made. This can assist in auditing user behavior, identifying potential security risks, and ensuring compliance with company policies.

```sql+postgres
select
  actor_email,
  actor_type,
  l.when,
  action ->> 'type' as action_type,
  action ->> 'result' as action_result,
  jsonb_pretty(new_value_json) as new_value,
  jsonb_pretty(old_value_json) as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  actor_email = 'user@domain.com';
```

```sql+sqlite
select
  actor_email,
  actor_type,
  l.when,
  json_extract(action, '$.type') as action_type,
  json_extract(action, '$.result') as action_result,
  new_value_json as new_value,
  old_value_json as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  actor_email = 'user@domain.com';
```

### Get all the activities performed on a particular resource
This query allows you to monitor and track all the activities carried out on a specific resource. It is beneficial for auditing purposes, such as identifying unauthorized changes or understanding user behavior.

```sql+postgres
select
  actor_email,
  actor_type,
  l.when,
  action ->> 'type' as action_type,
  action ->> 'result' as action_result,
  jsonb_pretty(new_value_json) as new_value,
  jsonb_pretty(old_value_json) as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  resource ->> 'id' = 'abcd13dcd91e9755b20ea5883fdd59ac';
```

```sql+sqlite
select
  actor_email,
  actor_type,
  l.when,
  json_extract(action, '$.type') as action_type,
  json_extract(action, '$.result') as action_result,
  new_value_json as new_value,
  old_value_json as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  json_extract(resource, '$.id') = 'abcd13dcd91e9755b20ea5883fdd59ac';
```

### Get all the activities performed on DNS records
Explore the actions performed on DNS records to gain insights into user activity and monitor changes. This can help in identifying unusual behavior or unauthorized modifications, enhancing the security and integrity of your DNS records.

```sql+postgres
select
  actor_email,
  actor_type,
  l.when,
  action ->> 'type' as action_type,
  action ->> 'result' as action_result,
  jsonb_pretty(new_value_json) as new_value,
  jsonb_pretty(old_value_json) as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  resource ->> 'type' = 'DNS_record';
```

```sql+sqlite
select
  actor_email,
  actor_type,
  l.when,
  json_extract(action, '$.type') as action_type,
  json_extract(action, '$.result') as action_result,
  new_value_json as new_value,
  old_value_json as old_value,
  owner_id
from
  cloudflare_user_audit_log l
where
  json_extract(resource, '$.type') = 'DNS_record';
```