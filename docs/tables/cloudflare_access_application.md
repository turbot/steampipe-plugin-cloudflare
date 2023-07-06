# Table: cloudflare_access_application

Access applications are used to restrict access to a whole application using an authorization gateway managed by Cloudflare.

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

### Get application count by account

```sql
select
  count(*),
  type
from
  cloudflare_access_application
group by
  account_id;
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
