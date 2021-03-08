# Table: cloudflare_rule_list

Rule Lists are a set of IP addresses or CIDR ranges that are configured on the account level. Once created, Rule Lists can be used in Firewall Rules across all zones within the same account.

## Examples

### List all Rule lists for the account

```sql
select
  *
from
  cloudflare_ip_list
```

### List all rule list items

```sql
select
  name,
  item ->> 'ip' as ip,
  item ->> 'comment' as comment
from
  cloudflare_rule_list,
  jsonb_array_elements(items) as item
```
