# Table: cloudflare_r2_object_data

Cloudflare R2 objects are stored in one or more Cloudflare R2 buckets, and each object can be up to 5 TB in size.

> Note: Using this table adds to cost to your monthly bill from Cloudflare. Optimizations have been put in place to minimize the impact as much as possible. Please refer to Cloudflare R2 Pricing to understand the cost implications.

## Examples

### Basic info

```sql
select
  key,
  etag,
  size
from
  cloudflare_r2_object
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01';
```

### List all objects with a fixed `prefix`

```sql
select
  key,
  etag,
  size,
  prefix
from
  cloudflare_r2_object
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and prefix = '/logs/2021/03/01/12';
```

### List all objects with a fixed `key`

```sql
select
  key,
  etag,
  size,
  prefix
from
  cloudflare_r2_object
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

### List all objects which were not modified in the last 3 months

```sql
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and prefix = 'static_assets'
  and last_modified < current_date - interval '3 months';
```
