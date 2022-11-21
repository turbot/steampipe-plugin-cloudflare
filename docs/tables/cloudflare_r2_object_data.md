# Table: cloudflare_r2_object_data

The `cloudflare_r2_object_data` table provides access to the data stored in a R2 object.

The data is serialized into a string if it contains valid UTF8 bytes, otherwise it is encoded into Base64, as defined in [RFC 4648](https://datatracker.ietf.org/doc/html/rfc4648).

To list objects, you must mention the `key` and the name of the container `bucket` which contains the objects.

> Note: Using this table adds to cost to your monthly bill from Cloudflare. Optimizations have been put in place to minimize the impact as much as possible. Please refer to [Cloudflare R2 Pricing](https://developers.cloudflare.com/r2/platform/pricing/) to understand the cost implications.

## Examples

### Basic info

```sql
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### Parse object data into `jsonb`

```sql
select
  key,
  bucket,
  data::jsonb
from
  cloudflare_r2_object_data
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### Process `jsonb` data in objects

```sql
select
  event ->> 'level' as level,
  event ->> 'severity' as severity,
  event ->> 'message' as event_message,
  event ->> 'data' as event_data,
  event ->> 'timestamp' as timestamp
from
  cloudflare_r2_object_data,
  jsonb_array_elements((data::jsonb) -> 'events') as event
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/auth_logs.json.gz'
  and event ->> 'level' = 'error';
```

### Get the raw binary `data` by converting back from `base64`

```sql
select
  decode(data, 'base64')
from
  cloudflare_r2_object_data
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'user_uploads'
  and key = 'avatar_9ac3097c-1e56-4108-b92e-226a3f4caeb8';
```
