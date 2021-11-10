# Table: cloudflare_worker_route

Routes are basic patterns that allow users to map a URL pattern to a Worker script to enable Workers to run on custom domains.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

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

### List worker routes not attached to any worker

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
