---
title: "Steampipe Table: cloudflare_user - Query Cloudflare Users using SQL"
description: "Allows users to query User data in Cloudflare, specifically the account details, email, status, and two-factor authentication details."
---

# Table: cloudflare_user - Query Cloudflare Users using SQL

Cloudflare is a global cloud network platform that offers various services to secure and enhance the performance of websites, applications, and other internet properties. It provides a range of services from content delivery network (CDN), website security, DDoS protection, to DNS services. A User in Cloudflare represents an individual with access to the Cloudflare dashboard and API.

## Table Usage Guide

The `cloudflare_user` table provides insights into User data within Cloudflare. As a Security Analyst, explore user-specific details through this table, including account details, email, status, and two-factor authentication details. Utilize it to uncover information about users, such as those with two-factor authentication enabled or disabled, and the verification of user status.

## Examples

### Query information about the user
Discover the details of your user profile on Cloudflare to better understand account status and settings. This can be useful for auditing purposes or for troubleshooting account-related issues.

```sql+postgres
select
  *
from
  cloudflare_user;
```

```sql+sqlite
select
  *
from
  cloudflare_user;
```

### Check if two factor authentication is enabled for the user
Explore whether two-factor authentication is activated for users, enhancing account security and reducing the risk of unauthorized access.

```sql+postgres
select
  id,
  email,
  two_factor_authentication_enabled
from
  cloudflare_user;
```

```sql+sqlite
select
  id,
  email,
  two_factor_authentication_enabled
from
  cloudflare_user;
```
