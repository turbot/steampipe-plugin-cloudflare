---
title: "Steampipe Table: cloudflare_account_role - Query Cloudflare Account Roles using SQL"
description: "Allows users to query Cloudflare Account Roles, providing detailed information about each role's permissions and associated metadata."
---

# Table: cloudflare_account_role - Query Cloudflare Account Roles using SQL

Cloudflare Account Roles are a part of Cloudflare's user management system that allows you to control who has access to your Cloudflare account and what they can do. Each role has a predefined set of permissions that determine what actions a user can take and what information they can view within your account. The roles system is designed to provide flexibility and control over your account's security.

## Table Usage Guide

The `cloudflare_account_role` table provides insights into the roles within Cloudflare's user management system. As a security analyst, you can explore role-specific details through this table, including permissions and associated metadata. Utilize it to uncover information about roles, such as those with specific permissions, the roles assigned to specific users, and to verify the security measures in place for your account.

## Examples

### Basic info
Explore which roles are associated with different accounts on Cloudflare, enabling you to manage and organize permissions effectively.

```sql+postgres
select
  name,
  id,
  account_id,
from
  cloudflare_account_role;
```

```sql+sqlite
select
  name,
  id,
  account_id
from
  cloudflare_account_role;
```