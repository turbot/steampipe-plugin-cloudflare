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
