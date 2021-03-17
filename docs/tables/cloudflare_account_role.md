# Table: cloudflare_account_role

Query information about the account roles, defines what permissions a Member of an Account has.

## Examples

### Basic info

```sql
select
  name,
  id,
  account_id,
from
  cloudflare_account_role;
```
