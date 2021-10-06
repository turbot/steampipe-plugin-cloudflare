# Table: cloudflare_worker_script

Routes are basic patterns used to enable or disable workers that match requests.

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
