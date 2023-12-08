---
title: "Steampipe Table: cloudflare_dns_record - Query Cloudflare DNS Records using SQL"
description: "Allows users to query DNS Records in Cloudflare, providing insights into the configuration details, zone information, and other metadata associated with each DNS record."
---

# Table: cloudflare_dns_record - Query Cloudflare DNS Records using SQL

Cloudflare DNS is a service that provides global, fast, and secure Domain Name System services. It is built on a network that is unified, distributed, and operates in real-time, providing users with reliable and highly available DNS services. It also offers DNSSEC to protect against forged DNS answers and other security threats.

## Table Usage Guide

The `cloudflare_dns_record` table provides insights into DNS records within Cloudflare. As a network administrator, you can explore record-specific details through this table, including the type of record, associated zone, and configuration settings. Utilize it to uncover information about DNS records, such as those with certain configurations, the zones they are associated with, and their current status.

**Important Notes**
- You must specify the `zone_id` in the `where` clause to query this table.

## Examples

### Query all DNS records for the zone
Explore all DNS records associated with a specific zone to understand its configuration and manage its settings effectively. This can be particularly useful in troubleshooting or optimizing network performance.

```sql+postgres
select
  *
from
  cloudflare_dns_record
where
  zone_id = 'ecfee56e04ffb0de172231a027abe23b';
```

```sql+sqlite
select
  *
from
  cloudflare_dns_record
where
  zone_id = 'ecfee56e04ffb0de172231a027abe23b';
```

### List MX records in priority order
Determine the areas in which MX records are prioritized within a specific zone in your Cloudflare DNS, providing a clear order of priority. This is useful for managing mail exchange servers and ensuring smooth email delivery.

```sql+postgres
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
  priority;
```

```sql+sqlite
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
  priority;
```

### List all records from each zone
Explore which DNS records belong to each zone in your Cloudflare account. This allows you to understand the distribution and organization of your DNS records, aiding in efficient management and troubleshooting.

```sql+postgres
select
  r.*,
  z.name as zone
from
  cloudflare_dns_record r,
  cloudflare_zone z
where
  r.zone_id = z.id;
```

```sql+sqlite
select
  r.*,
  z.name as zone
from
  cloudflare_dns_record r,
  cloudflare_zone z
where
  r.zone_id = z.id;
```