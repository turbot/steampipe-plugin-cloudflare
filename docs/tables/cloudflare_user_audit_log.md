# Table: cloudflare_user_audit_log

Information about the actions performed by users in the Cloudflare account.

## Examples

### Basic info

```sql
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

### Get all the users' activities in the last 10 days

```sql
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

### Get all the users' activities for a particular timeline

```sql
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

### Get all the activities of a particular user

```sql
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

### Get all the activities performed on a particular resource

```sql
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

### Get all the activities performed on DNS records

```sql
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