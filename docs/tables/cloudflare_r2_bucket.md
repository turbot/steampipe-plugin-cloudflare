# Table: cloudflare_r2_bucket

Cloudflare R2 Storage allows developers to store large amounts of unstructured data without the costly egress bandwidth fees associated with typical cloud storage services.

## Examples

### Basic info

```sql
select
  name,
  creation_date,
  region,
  account_id
from
  cloudflare_r2_bucket
where
  account_id = 'fb1696f453testaccount39e734f5f96e9';
```

### List buckets with default encryption disabled

```sql
select
  name,
  server_side_encryption_configuration
from
  cloudflare_r2_bucket
where
  server_side_encryption_configuration is null
  and account_id = 'fb1696f453testaccount39e734f5f96e9';
```
