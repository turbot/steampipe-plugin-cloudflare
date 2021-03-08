# Table: cloudflare_api_token

List all API tokens owned by the user.

**Warning**: This table is only available when using `email` and `api_key` for
authentication credentials. It does not work when using `api_token` for access.

## Examples

### Query all API tokens for this user

```sql
select
  *
from
  cloudflare_api_token
```

### List API tokens by age, oldest first

```sql
select
  id,
  name,
  status,
  issued_on,
  date_part('day', now() - issued_on) as age
from
  cloudflare_api_token
order by
  issued_on desc
```

### List API tokens expiring with the next 14 days

```sql
select
  id,
  name,
  status,
  expires_on
from
  cloudflare_api_token
where
  status is active
  and expires_on < current_date + interval '14 days'
```

### List all permissions granted to each API token

```sql
select
  name,
  policy ->> 'effect',
  policy -> 'resources',
  perm_group ->> 'name'
from
  cloudflare_api_token,
  jsonb_array_elements(policies) as policy,
  jsonb_array_elements(policy -> 'permission_groups') as perm_group
```
