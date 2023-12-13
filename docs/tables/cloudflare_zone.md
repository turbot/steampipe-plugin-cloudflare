---
title: "Steampipe Table: cloudflare_zone - Query Cloudflare Zones using SQL"
description: "Allows users to query Cloudflare Zones, providing insights into the DNS settings, SSL/TLS configurations, and associated metadata of each zone."
---

# Table: cloudflare_zone - Query Cloudflare Zones using SQL

A Cloudflare Zone represents a domain name that is registered with Cloudflare. It includes settings related to DNS, SSL/TLS, and other features that help protect and speed up your website. It is a crucial component in managing the performance and security of your web presence.

## Table Usage Guide

The `cloudflare_zone` table provides insights into zones within Cloudflare. As a network administrator, explore zone-specific details through this table, including DNS settings, SSL/TLS configurations, and associated metadata. Utilize it to uncover information about zones, such as their security level, development mode status, and the original DNS servers.

## Examples

### Query all zones for the user
Explore all zones associated with your user account on Cloudflare. This allows you to see a comprehensive overview of all your zones, useful for managing multiple domains or subdomains.

```sql+postgres
select
  *
from
  cloudflare_zone;
```

```sql+sqlite
select
  *
from
  cloudflare_zone;
```

### List all settings for the zone
Explore the various settings for a specific zone to gain insights into its configuration and values. This can aid in understanding the zone's current setup and potentially identifying areas for optimization or troubleshooting.

```sql+postgres
select
  name,
  setting.key,
  setting.value
from
  cloudflare_zone,
  jsonb_each_text(settings) as setting;
```

```sql+sqlite
select
  name,
  setting.key,
  setting.value
from
  cloudflare_zone,
  json_each(settings) as setting;
```

### Get details of the TLS 1.3 setting
Explore the configuration of your Cloudflare zones to understand the status of the TLS 1.3 setting. This can help ensure your zones are utilizing the latest security protocols.

```sql+postgres
select
  name,
  settings ->> 'tls_1_3'
from
  cloudflare_zone;
```

```sql+sqlite
select
  name,
  json_extract(settings, '$.tls_1_3')
from
  cloudflare_zone;
```

### List all permissions available to the user for this zone
Discover the segments that outline the range of permissions a user has in a certain zone, giving a comprehensive overview of their access rights. This is beneficial in maintaining security and ensuring appropriate access levels.

```sql+postgres
select
  name,
  perm
from
  cloudflare_zone,
  jsonb_array_elements_text(permissions) as perm;
```

```sql+sqlite
select
  name,
  perm.value
from
  cloudflare_zone,
  json_each(permissions) as perm;
```

### Check DNSSEC status for zones
Analyze the security status of your domain zones to ensure DNSSEC, a crucial internet security protocol, is properly enabled. This is essential for protecting your website from DNS spoofing and other DNS-related attacks.

```sql+postgres
select
  name,
  dnssec ->> 'status'
from
  cloudflare_zone;
```

```sql+sqlite
select
  name,
  json_extract(dnssec, '$.status')
from
  cloudflare_zone;
```
