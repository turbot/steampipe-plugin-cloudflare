# Table: cloudflare_access_group

An access group is a set of rules that can be configured once and then quickly applied across many Access applications.

## Examples

### Basic info

```sql
select
  name,
  id,
  created_at
from
  cloudflare_access_group;
```

### List access group rules

```sql
select
  name,
  id,
  jsonb_pretty(include) as include,
  jsonb_pretty(exclude) as exclude,
  jsonb_pretty(require) as require
from
  cloudflare_access_group;
```
