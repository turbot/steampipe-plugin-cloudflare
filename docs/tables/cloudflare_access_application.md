# Table: cloudflare_access_application

Access Applications are used to restrict access to a whole application using an authorisation gateway managed by Cloudflare.

**Note:** It's required that `account_id` is set in `~/.steampipe/config/cloudflare.spc` or through `CLOUDFLARE_ACCOUNT_ID` environment variable to access this table.

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
