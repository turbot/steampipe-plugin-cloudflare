# Table: cloudflare_load_balancer_monitor

A monitor issues health checks at regular intervals to evaluate the health of an origin pool.

## Examples

### Basic info

```sql
select
  id,
  type,
  path,
  timeout,
  retries,
  interval,
  port,
  expected_codes
from
  cloudflare_load_balancer_monitor;
```

### Get information of monitors attached to pool

```sql
select
  p.id as pool_id,
  p.name as pool_name,
  p.enabled as pool_enabled,
  p.notification_email,
  m.id as monitor_id,
  m.description monitor_description
from
  cloudflare_load_balancer_pool p,
  cloudflare_load_balancer_monitor as m
where
  p.monitor = m.id;
```
