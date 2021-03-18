# Table: cloudflare_page_rule

A page rule describes target patterns for requests and actions to perform on matching requests.

## Examples

### Basic info

```sql
select
  id,
  zone_id,
  status,
  priority
from
  cloudflare_page_rule;
```

### Get all page rules which are not enabled

```sql
select
  id,
  zone_id,
  status
from
  cloudflare_page_rule
where
  status = 'disabled';
```

### Get all page rules for which Always Online feature of the Cloudflare Caching app is not enabled

```sql
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
