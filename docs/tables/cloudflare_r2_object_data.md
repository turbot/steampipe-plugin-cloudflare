---
title: "Steampipe Table: cloudflare_r2_object_data - Query Cloudflare R2 Objects using SQL"
description: "Allows users to query Cloudflare R2 Objects, specifically object data, providing insights into the data stored in R2 storage."
---

# Table: cloudflare_r2_object_data - Query Cloudflare R2 Objects using SQL

Cloudflare R2 is a storage service that offers a simple, scalable, and cost-effective way to store and retrieve any amount of data at any time. It is designed to deliver 99.999999999% durability, and scale past trillions of objects worldwide. Customers use R2 for backups, restores, and to serve user-generated content.

## Table Usage Guide

The `cloudflare_r2_object_data` table provides insights into the objects stored in Cloudflare R2 storage. As a data analyst or a DevOps engineer, you can explore object-specific details through this table, including object metadata, storage class, and associated data. Utilize it to uncover information about objects, such as their size, last modified date, and the storage class they belong to.

**Important Notes**
- You must specify both the `key` and `bucket` in the `where` clause to query this table.
- Using this table adds to cost to your monthly bill from Cloudflare. Optimizations have been put in place to minimize the impact as much as possible. Please refer to [Cloudflare R2 Pricing](https://developers.cloudflare.com/r2/platform/pricing/) to understand the cost implications.

## Examples

### Basic info
Explore which types of content are stored in a specific Cloudflare account and bucket. This can be useful for understanding the structure and organization of your data, particularly for large-scale log management.

```sql+postgres
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

```sql+sqlite
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
Analyze the settings to understand the specific object data within a given Cloudflare account and bucket. This is particularly useful for exploring and understanding application log data.

```sql+postgres
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

```sql+sqlite
select
  key,
  bucket,
  json(data)
from
  cloudflare_r2_object_data
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### Process `jsonb` data in objects
Determine the areas in which errors occur in your application by analyzing event data from your Cloudflare logs. This allows you to pinpoint specific instances of error levels, enhancing your ability to troubleshoot and improve your application's performance.

```sql+postgres
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

```sql+sqlite
select
  json_extract(event.value, '$.level') as level,
  json_extract(event.value, '$.severity') as severity,
  json_extract(event.value, '$.message') as event_message,
  json_extract(event.value, '$.data') as event_data,
  json_extract(event.value, '$.timestamp') as timestamp
from
  cloudflare_r2_object_data,
  json_each(json_extract(data, '$.events')) as event
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/auth_logs.json.gz'
  and json_extract(event.value, '$.level') = 'error';
```

### Get the raw binary `data` by converting back from `base64`
Discover the segments that enable the extraction of raw binary data from a specific user's uploaded files in a Cloudflare account. This might be used to analyze or manipulate the file data directly, bypassing the need for base64 encoding.

```sql+postgres
select
  decode(data, 'base64')
from
  cloudflare_r2_object_data
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'user_uploads'
  and key = 'avatar_9ac3097c-1e56-4108-b92e-226a3f4caeb8';
```

```sql+sqlite
select
  decode(data, 'base64')
from
  cloudflare_r2_object_data
where
  account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'user_uploads'
  and key = 'avatar_9ac3097c-1e56-4108-b92e-226a3f4caeb8';
```

### List the object data of those objects that are encrypted with SSE KMS key
Explore which objects within a specific account and bucket are encrypted using a KMS key. This is particularly useful for identifying and managing sensitive data that requires enhanced security measures.

```sql+postgres
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  sse_kms_key_id is not null
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

```sql+sqlite
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  sse_kms_key_id is not null
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### List the object data of those objects that are expiring in the next 7 days
Determine the details of certain objects set to expire within the next week. This could be particularly useful for managing and prioritizing updates or renewals for those objects, especially within a large dataset.

```sql+postgres
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  expires >= now() + interval '7' day
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

```sql+sqlite
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  expires >= datetime('now', '+7 days')
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### List the object data of those objects that are 'Delete Marker'
Determine the areas in which certain objects are marked for deletion within a specific account and bucket. This is particularly useful in identifying and managing potential data removals.

```sql+postgres
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  delete_marker 
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

```sql+sqlite
select
  key,
  bucket,
  content_type
from
  cloudflare_r2_object_data
where
  delete_marker
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```