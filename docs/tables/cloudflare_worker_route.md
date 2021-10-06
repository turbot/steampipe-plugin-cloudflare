# Table: cloudflare_worker_script

Routes are basic patterns that allow users to map a URL pattern to a Worker script to enable Workers to run on custom domains.

**Note:** It's required that `account_id` is set in `~/.steampipe/config/cloudflare.spc` or through `CLOUDFLARE_ACCOUNT_ID` environment variable to access this table.

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
