# Table: cloudflare_access_policy

Access Policies are used in conjunction with Access Applications to restrict access to a particular resource.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

## Examples

### Basic info

```sql
select
  name,
  id,
  application_id,
  application_name,
  decision,
  precedence
from
  cloudflare_access_policy;
```

### List policies that require justifcation for accessing resources

```sql
select
  name,
  id,
  application_name,
  decision,
  precedence
from
  cloudflare_access_policy
where
  purpose_justification_required;
```
