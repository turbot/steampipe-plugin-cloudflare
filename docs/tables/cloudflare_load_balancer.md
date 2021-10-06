# Table: cloudflare_load_balancer

Cloudflare Load balancers allows to distribute traffic across servers, which reduces server strain and latency and improves the experience for end users.

## Examples

### Basic info

```sql
select
  name,
  id,
  zone_name,
  zone_id,
  created_on,
  ttl,
  steering_policy
from
  cloudflare_load_balancer;
```

### List proxied load balancers

```sql
select
  name,
  zone_id,
  ttl
from
  cloudflare_load_balancer
where
  proxied;
```

### Get session_affinity details for load balancers

```sql
select
  name,
  zone_name,
  session_affinity,
  session_affinity_ttl,
  jsonb_pretty(session_affinity_attributes) as session_affinity_attributes
from
  cloudflare_load_balancer;
```
