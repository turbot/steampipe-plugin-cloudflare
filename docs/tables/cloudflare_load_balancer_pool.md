# Table: cloudflare_load_balancer_pool

A pool is a group of origin servers, with each origin identified by its IP address or hostname.

## Examples

### Basic info

```sql
select
  name,
  id,
  enabled,
  description,
  created_on,
  jsonb_pretty(origins) as origins
from
  cloudflare_load_balancer_pool;
```

### List active pools

```sql
select
  name,
  id,
  monitor
from
  cloudflare_load_balancer_pool
where
  enabled;
```
