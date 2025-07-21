---
title: "Steampipe Table: cloudflare_zone_setting - Query Cloudflare Zone Settings using SQL"
description: "Allows users to query individual zone settings in Cloudflare, providing detailed information about each setting's value and configuration."
---

# Table: cloudflare_zone_setting - Query Cloudflare Zone Settings using SQL

Cloudflare Zone Settings control various features and behaviors for a zone, such as security level, cache settings, SSL configuration, and performance optimizations. Each setting has a specific ID and value that determines how Cloudflare handles requests for that zone.

## Table Usage Guide

The `cloudflare_zone_setting` table provides insights into individual zone settings within Cloudflare. As a DevOps engineer or security analyst, explore setting-specific details through this table, including current values, editability, and modification timestamps. Utilize it to audit zone configurations, understand security policies, and manage feature settings across your zones.

**Important Notes:**
- By default this table fetch all settings across all zones.
- For optimal performance and to reduce query time, always specify `zone_id` and/or `id` (setting ID) in your WHERE clause
- Possible values for `id` are:
    - `0rtt`
    - `advanced_ddos`
    - `always_online`
    - `always_use_https`
    - `automatic_https_rewrites`
    - `brotli`
    - `browser_cache_ttl`
    - `browser_check`
    - `cache_level`
    - `challenge_ttl`
    - `ciphers`
    - `cname_flattening`
    - `development_mode`
    - `early_hints`
    - `edge_cache_ttl`
    - `email_obfuscation`
    - `h2_prioritization`
    - `hotlink_protection`
    - `http2`
    - `http3`
    - `image_resizing`
    - `ip_geolocation`
    - `ipv6`
    - `max_upload`
    - `min_tls_version`
    - `mirage`
    - `nel`
    - `opportunistic_encryption`
    - `opportunistic_onion`
    - `orange_to_orange`
    - `origin_error_page_pass_thru`
    - `origin_h2_max_streams`
    - `origin_max_http_version`
    - `polish`
    - `prefetch_preload`
    - `privacy_pass`
    - `proxy_read_timeout`
    - `pseudo_ipv4`
    - `replace_insecure_js`
    - `response_buffering`
    - `rocket_loader`
    - `automatic_platform_optimization`
    - `security_header`
    - `security_level`
    - `server_side_exclude`
    - `sha1_support`
    - `sort_query_string_for_cache`
    - `ssl`
    - `ssl_recommender`
    - `tls_1_2_only`
    - `tls_1_3`
    - `tls_client_auth`
    - `true_client_ip_header`
    - `waf`
    - `webp`
    - `websockets`

## Examples

### List all settings for a specific zone
Explore all available settings for a particular zone to understand its current configuration. Always specify `zone_id` for better performance.

```sql+postgres
select
  id,
  value,
  editable,
  enabled,
  modified_on
from
  cloudflare_zone_setting
where
  zone_id = '41b12a8232fe413913ddef4714c0f19b';
```

```sql+sqlite
select
  id,
  value,
  editable,
  enabled,
  modified_on
from
  cloudflare_zone_setting
where
  zone_id = '41b12a8232fe413913ddef4714c0f19b';
```

### Get a specific setting for a specific zone
Retrieve a single setting for a zone by specifying both `zone_id` and setting `id`. This is the most efficient query pattern.

```sql+postgres
select
  id,
  value,
  editable,
  modified_on
from
  cloudflare_zone_setting
where
  zone_id = '41b12a8232fe413913ddef4714c0f19b'
  and id = 'ssl';
```

```sql+sqlite
select
  id,
  value,
  editable,
  modified_on
from
  cloudflare_zone_setting
where
  zone_id = '41b12a8232fe413913ddef4714c0f19b'
  and id = 'ssl';
```

### Check SSL/TLS settings across all zones with zone details
Examine SSL and TLS configuration settings across all your zones, joining with the zone table to get zone names.

```sql+postgres
select
  z.name as zone_name,
  zs.id as setting_id,
  zs.value,
  zs.editable,
  zs.modified_on
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.id in ('ssl', 'tls_1_3', 'min_tls_version', 'always_use_https')
order by
  z.name, zs.id;
```

```sql+sqlite
select
  z.name as zone_name,
  zs.id as setting_id,
  zs.value,
  zs.editable,
  zs.modified_on
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.id in ('ssl', 'tls_1_3', 'min_tls_version', 'always_use_https')
order by
  z.name, zs.id;
```

### Find zones with specific security settings
Identify zones that have specific security features enabled or disabled, including zone metadata.

```sql+postgres
select
  z.name as zone_name,
  z.status as zone_status,
  zs.id as setting_id,
  zs.value,
  zs.editable
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.id in ('security_level', 'browser_check', 'challenge_ttl')
  and zs.value != 'off'
order by
  z.name;
```

```sql+sqlite
select
  z.name as zone_name,
  z.status as zone_status,
  zs.id as setting_id,
  zs.value,
  zs.editable
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.id in ('security_level', 'browser_check', 'challenge_ttl')
  and zs.value != 'off'
order by
  z.name;
```

### Query specific setting across all zones with zone info
Get the value of a particular setting across all your zones, including zone details.

```sql+postgres
select
  z.name as zone_name,
  z.type as zone_type,
  z.status as zone_status,
  zs.value,
  zs.editable,
  zs.modified_on
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.id = 'development_mode'
order by
  z.name;
```

```sql+sqlite
select
  z.name as zone_name,
  z.type as zone_type,
  z.status as zone_status,
  zs.value,
  zs.editable,
  zs.modified_on
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.id = 'development_mode'
order by
  z.name;
```

### List editable vs non-editable settings summary
Understand which settings can be modified for your zones.

```sql+postgres
select
  id as setting_id,
  editable,
  count(*) as zone_count,
  count(case when enabled = true then 1 end) as enabled_count
from
  cloudflare_zone_setting
group by
  id, editable
order by
  id;
```

```sql+sqlite
select
  id as setting_id,
  editable,
  count(*) as zone_count,
  count(case when enabled = 1 then 1 end) as enabled_count
from
  cloudflare_zone_setting
group by
  id, editable
order by
  id;
```

### Find recently modified settings across zones
Identify settings that have been recently modified, useful for auditing configuration changes.

```sql+postgres
select
  z.name as zone_name,
  zs.id as setting_id,
  zs.value,
  zs.modified_on
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.modified_on >= now() - interval '7 days'
order by
  zs.modified_on desc;
```

```sql+sqlite
select
  z.name as zone_name,
  zs.id as setting_id,
  zs.value,
  zs.modified_on
from
  cloudflare_zone_setting zs
  join cloudflare_zone z on zs.zone_id = z.id
where
  zs.modified_on >= datetime('now', '-7 days')
order by
  zs.modified_on desc;
```

### Compare settings between zones
Compare specific settings between different zones to ensure consistency.

```sql+postgres
select
  zs1.id as setting_id,
  z1.name as zone1_name,
  zs1.value as zone1_value,
  z2.name as zone2_name,
  zs2.value as zone2_value,
  case
    when zs1.value = zs2.value then 'Match'
    else 'Different'
  end as comparison
from
  cloudflare_zone_setting zs1
  join cloudflare_zone z1 on zs1.zone_id = z1.id
  join cloudflare_zone_setting zs2 on zs1.id = zs2.id
  join cloudflare_zone z2 on zs2.zone_id = z2.id
where
  z1.name = 'pleasehelpme.com'
  and z2.name = 'myactualdomain.com'
  and zs1.id in ('ssl', 'security_level', 'cache_level')
order by
  zs1.id;
```

```sql+sqlite
select
  zs1.id as setting_id,
  z1.name as zone1_name,
  zs1.value as zone1_value,
  z2.name as zone2_name,
  zs2.value as zone2_value,
  case
    when zs1.value = zs2.value then 'Match'
    else 'Different'
  end as comparison
from
  cloudflare_zone_setting zs1
  join cloudflare_zone z1 on zs1.zone_id = z1.id
  join cloudflare_zone_setting zs2 on zs1.id = zs2.id
  join cloudflare_zone z2 on zs2.zone_id = z2.id
where
  z1.name = 'pleasehelpme.com'
  and z2.name = 'myactualdomain.com'
  and zs1.id in ('ssl', 'security_level', 'cache_level')
order by
  zs1.id;
```