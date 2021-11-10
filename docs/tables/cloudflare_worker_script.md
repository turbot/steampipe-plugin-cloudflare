# Table: cloudflare_worker_script

A Worker script is a single script that is executed on matching routes in the Cloudflare edge.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

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
