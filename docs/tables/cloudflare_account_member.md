---
title: "Steampipe Table: cloudflare_account_member - Query Cloudflare Account Members using SQL"
description: "Allows users to query Cloudflare Account Members, providing detailed information about members associated with a Cloudflare account."
---

# Table: cloudflare_account_member - Query Cloudflare Account Members using SQL

A Cloudflare Account Member represents an individual user in the context of a Cloudflare account. Account Members have various roles and permissions that determine what actions they can perform and what information they can access within the account. This includes access to various Cloudflare services, data, and configurations.

## Table Usage Guide

The `cloudflare_account_member` table provides insights into individual users within a Cloudflare account. As a system administrator, explore member-specific details through this table, including roles, permissions, and associated metadata. Utilize it to uncover information about members, such as their access levels, the services they can manage, and the configurations they can modify.

## Examples

### List all members in an account
Explore which members are associated with your account to manage user access and permissions more effectively. This is useful for maintaining security and ensuring the right people have access to the right resources.

```sql+postgres
select
  title,
  account_id,
  user_email
from
  cloudflare_account_member;
```

```sql+sqlite
select
  title,
  account_id,
  user_email
from
  cloudflare_account_member;
```

### List of members with Administrator access
Discover the segments that have members with Administrator access. This query can be used to identify potential security risks by highlighting accounts where users have been granted high-level permissions.

```sql+postgres
select
  title,
  account_id,
  user_email,
  attached_roles ->> 'name' as attached_role_name
from
  cloudflare_account_member,
  jsonb_array_elements(roles) as attached_roles
where
  attached_roles ->> 'name' = 'Administrator';
```

```sql+sqlite
select
  title,
  account_id,
  user_email,
  json_extract(attached_roles.value, '$.name') as attached_role_name
from
  cloudflare_account_member,
  json_each(roles) as attached_roles
where
  json_extract(attached_roles.value, '$.name') = 'Administrator';
```

### List of members yet to accept the join request
Discover the segments that include members who have not yet accepted their requests to join your Cloudflare account. This is useful for tracking pending invitations and managing team access to your account.

```sql+postgres
select
  title,
  account_id,
  user_email,
  status
from
  cloudflare_account_member
where
  status = 'pending';
```

```sql+sqlite
select
  title,
  account_id,
  user_email,
  status
from
  cloudflare_account_member
where
  status = 'pending';
```