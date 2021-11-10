# Table: cloudflare_workers_kv_namespace

Workers KV is a global, low-latency, key-value data store. It supports exceptionally high read volumes with low-latency, making it possible to build highly dynamic APIs and websites which respond as quickly as a cached static file would.
A Namespace is a collection of key-value pairs stored in Workers KV.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

## Examples

### Basic info

```sql
select
  id,
  title
from
  cloudflare_workers_kv_namespace;
```
