---
title: "Steampipe Table: cloudflare_r2_bucket - Query Cloudflare R2 Buckets using SQL"
description: "Allows users to query Cloudflare R2 Buckets, specifically providing details about each bucket's configuration, status, and usage."
---

# Table: cloudflare_r2_bucket - Query Cloudflare R2 Buckets using SQL

Cloudflare R2 is a cloud storage solution that offers high-performance, scalable, and secure object storage. It allows developers to store and retrieve any amount of data at any time from anywhere on the web. R2 buckets are the fundamental containers in Cloudflare R2 for data storage.

## Table Usage Guide

The `cloudflare_r2_bucket` table provides insights into R2 buckets within Cloudflare. As a cloud engineer or developer, you can explore bucket-specific details through this table, including the bucket's configuration, status, and usage. Utilize it to manage and monitor your Cloudflare R2 storage, ensuring optimal performance and security.

## Examples

### Basic info
Explore which Cloudflare resources were created in a specific region within a given account. This could be useful to manage resource allocation and understand usage patterns across different regions.

```sql+postgres
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

```sql+sqlite
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
Explore which Cloudflare R2 buckets lack default encryption, which could potentially expose sensitive data. This query is particularly useful for identifying security vulnerabilities in your storage configuration.

```sql+postgres
select
  name,
  server_side_encryption_configuration
from
  cloudflare_r2_bucket
where
  server_side_encryption_configuration is null
  and account_id = 'fb1696f453testaccount39e734f5f96e9';
```

```sql+sqlite
select
  name,
  server_side_encryption_configuration
from
  cloudflare_r2_bucket
where
  server_side_encryption_configuration is null
  and account_id = 'fb1696f453testaccount39e734f5f96e9';
```

### List buckets created in the last 30 days
Explore which buckets have been created in the last 30 days for a specific account. This allows you to keep track of recent additions and understand the level of server-side encryption configuration for each of these new buckets.

```sql+postgres
select
  name,
  server_side_encryption_configuration
from
  cloudflare_r2_bucket
where
  creation_date >= now() - interval '30' day
  and account_id = 'fb1696f453testaccount39e734f5f96e9';
```

```sql+sqlite
select
  name,
  server_side_encryption_configuration
from
  cloudflare_r2_bucket
where
  creation_date >= datetime('now', '-30 days')
  and account_id = 'fb1696f453testaccount39e734f5f96e9';
```