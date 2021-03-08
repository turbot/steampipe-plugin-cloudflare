# Table: cloudflare_zone

Zone is the basic resource for working with Cloudflare and is roughly equivalent to a domain name that the user purchases.

## Examples

### Query all zones for the user

```sql
select
  *
from
  cloudflare_zone
```

### List all settings for the zone

```sql
select
  name,
  setting.key,
  setting.value
from
  cloudflare_zone,
  jsonb_each_text(settings) as setting
```

### Get details of the TLS 1.3 setting

```sql
select
  name,
  settings ->> 'tls_1_3'
from
  cloudflare_zone
```

### List all permissions available to the user for this zone

```sql
select
  name,
  perm
from
  cloudflare_zone,
  jsonb_array_elements_text(permissions) as perm
```

### Check DNSSEC status for zones

```sql
select
  name,
  dnssec ->> 'status'
from
  cloudflare_zone
```
