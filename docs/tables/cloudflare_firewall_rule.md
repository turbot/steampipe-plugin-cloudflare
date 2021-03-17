# Table: cloudflare_firewall_rule

Query information about firewall rules. It uses filter expressions for more control over how traffic is matched to the rule. A filter expression permits selecting traffic by multiple criteria allowing greater freedom in rule creation.

## Examples

### Basic info

```sql
select
  id,
  zone_id,
  created_on
from
  cloudflare_firewall_rule;
```

### Query all firewall rules which are paused

```sql
select
  id,
  zone_id,
  paused
from
  cloudflare_firewall_rule
where
  paused;
```

### Get the rule that blocks requests based on IP reputation

```sql
select
  id,
  zone_id,
  filter,
  action
from
  cloudflare_firewall_rule
where
  action = 'block'
  and filter ->> 'expression' = '(cf.threat_score gt 1)';
```
