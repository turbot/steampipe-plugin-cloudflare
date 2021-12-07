# Table: cloudflare_access_group

A group is a set of rules that can be configured once and then quickly applied across many Access applications.
Access group allows to define a set of users to which an application policy can be applied.

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