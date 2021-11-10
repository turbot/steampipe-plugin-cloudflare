# Table: cloudflare_access_application

Access Applications are used to restrict access to a whole application using an authorisation gateway managed by Cloudflare.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

## Examples

### Basic info

```sql
select
  name,
  id,
  domain,
  created_at
from
  cloudflare_access_application;
```

### Get applications count by type

```sql
select
  count(*),
  type
from
  cloudflare_access_application
group by
  type;
```

### List applications with binding cookie enabled for increased security

```sql
select
  name,
  id,
  domain
from
  cloudflare_access_application
where
  enable_binding_cookie;
```
