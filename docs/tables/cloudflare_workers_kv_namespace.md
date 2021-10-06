# Table: cloudflare_workers_kv_namespace

Workers KV is a global, low-latency, key-value data store. It supports exceptionally high read volumes with low-latency, making it possible to build highly dynamic APIs and websites which respond as quickly as a cached static file would.
A Namespace is a collection of key-value pairs stored in Workers KV.

**Note:** It's required that `account_id` is set in `~/.steampipe/config/cloudflare.spc` or through `CLOUDFLARE_ACCOUNT_ID` environment variable to access this table.

## Examples

### Basic info

```sql
select
  id,
  title
from
  cloudflare_workers_kv_namespace;
```
