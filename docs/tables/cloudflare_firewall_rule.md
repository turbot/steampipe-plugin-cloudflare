# Table: cloudflare_firewall_rule

Firewall rules use filter expressions to control what traffic is allowed. A filter expression permits selecting traffic by multiple criteria allowing greater freedom in rule creation.

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

### List paused firewall rules

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

### List firewall rules that block requests based on IP reputation

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
