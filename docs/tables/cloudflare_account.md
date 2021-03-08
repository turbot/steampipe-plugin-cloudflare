# Table: cloudflare_account

Query information about your current account.

## Examples

### Query all accounts the user has access to

```sql
select
  *
from
  cloudflare_account
```

### Check if two factor authentication is enforced for accounts

```sql
select
  name,
  settings -> 'enforce_twofactor' as enforce_mfa
from
  cloudflare_account
```
