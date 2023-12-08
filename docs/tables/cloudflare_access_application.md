---
title: "Steampipe Table: cloudflare_access_application - Query Cloudflare Access Applications using SQL"
description: "Allows users to query Cloudflare Access Applications, specifically the application's ID, name, domain, session duration, and other related details, providing insights into the applications' configurations and settings."
---

# Table: cloudflare_access_application - Query Cloudflare Access Applications using SQL

Cloudflare Access is a cloud-based security service that protects your internal applications without a VPN. It operates on a per-user basis, giving you control over who can access your internal applications, and from where. It also offers features like Single Sign-On (SSO) and Multi-Factor Authentication (MFA) for added security.

## Table Usage Guide

The `cloudflare_access_application` table provides insights into Access Applications within Cloudflare Access. As a security engineer, explore application-specific details through this table, including the application's domain, session duration, and access policies. Utilize it to uncover information about applications, such as their configuration, settings, and the security measures in place.

## Examples

### Basic info
Explore which Cloudflare access applications have been created, along with their respective names, IDs and domains. This can be beneficial in managing access control and understanding the distribution of applications across different domains.

```sql+postgres
select
  name,
  id,
  domain,
  created_at
from
  cloudflare_access_application;
```

```sql+sqlite
select
  name,
  id,
  domain,
  created_at
from
  cloudflare_access_application;
```

### Get application count by account
Gain insights into the distribution of applications across different accounts. This query is useful for understanding the usage patterns and managing resources efficiently.

```sql+postgres
select
  count(*),
  type
from
  cloudflare_access_application
group by
  account_id;
```

```sql+sqlite
select
  count(*),
  type
from
  cloudflare_access_application
group by
  account_id;
```

### List applications with binding cookie enabled for increased security
Analyze the settings to understand which applications have the binding cookie enabled for increased security. This is useful for identifying potential vulnerabilities and ensuring optimal security configurations.

```sql+postgres
select
  name,
  id,
  domain
from
  cloudflare_access_application
where
  enable_binding_cookie;
```

```sql+sqlite
select
  name,
  id,
  domain
from
  cloudflare_access_application
where
  enable_binding_cookie = 1;
```