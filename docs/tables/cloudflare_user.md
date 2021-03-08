# Table: cloudflare_user

Information about the current user making the request.

## Examples

### Query information about the user

```sql
select
  *
from
  cloudflare_user
```

### Check if two factor authentication is enabled for the user

```sql
select
  id,
  email,
  two_factor_authentication_enabled
from
  cloudflare_user
```
