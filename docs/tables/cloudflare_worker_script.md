# Table: cloudflare_worker_script

A Worker script is a single script that is executed on matching routes in the Cloudflare edge.

**Note:** It's required that `account_id` is set in `~/.steampipe/config/cloudflare.spc` or through `CLOUDFLARE_ACCOUNT_ID` environment variable to access this table.

## Examples

### Basic info

```sql
select
  id,
  etag,
  created_on
from
  cloudflare_worker_script;
```
