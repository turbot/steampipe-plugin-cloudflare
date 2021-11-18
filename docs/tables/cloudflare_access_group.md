# Table: cloudflare_access_group

A group is a set of rules that can be configured once and then quickly applied across many Access applications.
Access group allows to define a set of users to which an application policy can be applied.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

**Warning**: If `account_id` is missing in the connection configuration. Query to this table will error out with message: `Error: HTTP status 400: Could not route to /accounts/access/groups, perhaps your object identifier is invalid? (7003), No route for that URI (7000)`

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
