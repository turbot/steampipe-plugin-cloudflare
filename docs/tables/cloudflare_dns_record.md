# Table: cloudflare_dns_record

List all DNS records associated with a zone.

Note: A `zone_id` must be provided in all queries to this table.

## Examples

### Query all DNS records for the zone

```sql
select
  *
from
  cloudflare_dns_record
where
  zone_id = 'ecfee56e04ffb0de172231a027abe23b'
```

### List MX records in priority order

```sql
select
  name,
  type,
  priority
from
  cloudflare_dns_record
where
  zone_id = 'ecfee56e04ffb0de172231a027abe23b'
  and type = 'MX'
order by
  priority
```
