---
title: "Steampipe Table: cloudflare_access_policy - Query Cloudflare Access Policies using SQL"
description: "Allows users to query Cloudflare Access Policies, specifically the policy details, providing insights into access control and potential security issues."
---

# Table: cloudflare_access_policy - Query Cloudflare Access Policies using SQL

Cloudflare Access Policy is a feature within Cloudflare that allows you to control who can access your application based on user identity, IP address, or other criteria. It provides a centralized way to set up and manage access policies for various Cloudflare resources, including web applications, databases, and more. Cloudflare Access Policy helps you secure your applications and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `cloudflare_access_policy` table provides insights into Access Policies within Cloudflare. As a Security Analyst, explore policy-specific details through this table, including permissions, IP addresses, and associated metadata. Utilize it to uncover information about policies, such as those with specific permissions, the IP addresses associated with policies, and the verification of access conditions.

## Examples

### Basic info
Explore which access policies are in place within your Cloudflare application. This is useful for assessing the precedence of these policies and understanding the decision-making process of your application's security.

```sql+postgres
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

```sql+sqlite
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

### List policies that require justification for accessing resources
Explore which policies necessitate a justification when accessing resources. This can be useful in enhancing security measures by identifying areas where additional user accountability is needed.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  application_name,
  decision,
  precedence
from
  cloudflare_access_policy
where
  purpose_justification_required = 1;
```