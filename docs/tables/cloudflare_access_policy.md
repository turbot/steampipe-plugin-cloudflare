# Table: cloudflare_access_policy

Access Policies are used in conjunction with Access Applications to restrict access to a particular resource.

**Note:** An account ID must be set in the connection configuration's `account_id` argument or through the `CLOUDFLARE_ACCOUNT_ID` environment variable to query this table.

**Warning**: If `account_id` is missing in the connection configuration. Query to this table will error out with message: `Error: HTTP status 400: Could not route to /accounts/access/apps, perhaps your object identifier is invalid? (7003), No route for that URI (7000)`

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
