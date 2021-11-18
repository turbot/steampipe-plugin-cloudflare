# Table: cloudflare_worker_route

Routes are basic patterns that allow users to map a URL pattern to a Worker script to enable Workers to run on custom domains.

## Examples

### Basic info

```sql
select
  id,
  zone_name,
  pattern,
  script
from
  cloudflare_worker_route;
```

### List idle worker routes (i.e. not attached to any worker)

```sql
select
  id,
  zone_name,
  pattern
from
  cloudflare.cloudflare_worker_route
where
  script = '';
```
