---
title: "Steampipe Table: cloudflare_access_group - Query Cloudflare Access Groups using SQL"
description: "Allows users to query Access Groups in Cloudflare, specifically the access group settings, providing insights into access controls and security configurations."
---

# Table: cloudflare_access_group - Query Cloudflare Access Groups using SQL

Cloudflare Access Groups is a feature within Cloudflare that allows you to manage and control access to your applications and services. It provides a way to set up and manage groups of users who have access to specific resources, based on predefined conditions. Cloudflare Access Groups help you maintain the security and integrity of your resources by ensuring only authorized users can access them.

## Table Usage Guide

The `cloudflare_access_group` table provides insights into Access Groups within Cloudflare. As a security analyst, explore group-specific details through this table, including group names, user emails, and associated metadata. Utilize it to uncover information about groups, such as those with specific access permissions, the users associated with each group, and the verification of access controls.


## Examples

### Basic info
Determine the areas in which Cloudflare access groups were established by examining their creation dates. This can help in understanding the timeline of security group deployment and aid in managing access control.

```sql+postgres
select
  name,
  id,
  created_at
from
  cloudflare_access_group;
```

```sql+sqlite
select
  name,
  id,
  created_at
from
  cloudflare_access_group;
```

### List access group rules
Analyze the settings to understand the rules of your access groups. This can help you pinpoint specific locations where access is granted or denied, providing a comprehensive view of your security configurations.

```sql+postgres
select
  name,
  id,
  jsonb_pretty(include) as include,
  jsonb_pretty(exclude) as exclude,
  jsonb_pretty(require) as require
from
  cloudflare_access_group;
```

```sql+sqlite
select
  name,
  id,
  include,
  exclude,
  require
from
  cloudflare_access_group;
```