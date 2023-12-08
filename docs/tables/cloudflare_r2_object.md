---
title: "Steampipe Table: cloudflare_r2_object - Query Cloudflare R2 Objects using SQL"
description: "Allows users to query R2 Objects in Cloudflare, specifically providing insights into object details such as object ID, object type, and object content."
---

# Table: cloudflare_r2_object - Query Cloudflare R2 Objects using SQL

Cloudflare R2 is a storage service that provides fast, reliable, and cost-effective object storage. It is designed to handle data from any source, making it a versatile solution for storing, retrieving, and managing data. With R2, users can store any amount of data and access it anywhere at any time.

## Table Usage Guide

The `cloudflare_r2_object` table provides insights into the R2 Objects within Cloudflare. As a DevOps engineer, explore object-specific details through this table, such as the object ID, object type, and object content. Utilize it to manage and monitor your Cloudflare R2 Objects, ensuring optimal data storage and retrieval.

**Important Notes**
- Using this table adds to cost to your monthly bill from Cloudflare. Optimizations have been put in place to minimize the impact as much as possible. Please refer to Cloudflare R2 Pricing to understand the cost implications.

## Examples

### Basic info
Explore the specific details of an account's stored objects within a particular Cloudflare bucket to understand their size and changes. This can be useful for assessing storage usage and tracking modifications to objects within a specific time period.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have a specific prefix within your Cloudflare account. This is beneficial for organizing and locating specific sets of data, such as logs from a certain date.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that contain specific objects within a given account and bucket, which can be useful for pinpointing where certain data is stored or identifying patterns in data storage.

```sql+postgres
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

```sql+sqlite
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
Explore the objects within a specific account and bucket that have remained unmodified over the past three months. This query can be used to identify stagnant or unused data, aiding in efficient data management and potential cleanup efforts.

```sql+postgres
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

```sql+sqlite
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
  and last_modified < date('now','-3 months');
```

### List objects created in the last 30 days
Explore objects that have been created in the past month within a specific account and bucket. This can be useful for monitoring recent activity or changes in your Cloudflare data.

```sql+postgres
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  creation_date >= now() - interval '30' day
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

```sql+sqlite
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  creation_date >= datetime('now', '-30 day')
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

### List objects that return truncated results
Determine the instances where certain objects are yielding incomplete results within a specific Cloudflare account and bucket, allowing you to identify and address potential issues with data integrity or retrieval.

```sql+postgres
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  is_truncated
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

```sql+sqlite
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  is_truncated
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

### List objects that have bucket key enabled
Explore which objects within a specific Cloudflare account have the bucket key enabled. This is particularly useful for identifying potential security risks or for auditing purposes.

```sql+postgres
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  bucket_key_enabled
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

```sql+sqlite
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  cloudflare_r2_object
where
  bucket_key_enabled = 1
  and account_id = 'fb1696f453testaccount39e734f5f96e9'
  and bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```