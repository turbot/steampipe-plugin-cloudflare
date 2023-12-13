---
title: "Steampipe Table: cloudflare_load_balancer - Query Cloudflare Load Balancers using SQL"
description: "Allows users to query Cloudflare Load Balancers, providing insights into the load balancing configurations and health status of the various resources."
---

# Table: cloudflare_load_balancer - Query Cloudflare Load Balancers using SQL

Cloudflare Load Balancers are part of Cloudflare's performance and reliability services. They distribute network or application traffic across many resources, effectively reducing the burden on any single resource and ensuring optimal resource utilization. They also provide health checks and failover support, enhancing the reliability of your services.

## Table Usage Guide

The `cloudflare_load_balancer` table provides insights into the load balancing configurations and health status within Cloudflare. As a network administrator or a DevOps engineer, you can explore load balancer-specific details through this table, including the traffic distribution, failover configurations, and associated metadata. Utilize it to uncover information about your load balancers, such as their health status, the resources they are balancing, and their failover configurations.

## Examples

### Basic info
Explore the setup of your load balancers in Cloudflare to understand their configuration and creation details, which can help in assessing their performance and identifying potential areas for optimization.

```sql+postgres
select
  name,
  id,
  zone_name,
  zone_id,
  created_on,
  ttl,
  steering_policy
from
  cloudflare_load_balancer;
```

```sql+sqlite
select
  name,
  id,
  zone_name,
  zone_id,
  created_on,
  ttl,
  steering_policy
from
  cloudflare_load_balancer;
```

### List proxied load balancers
Explore which load balancers are set to be proxied. This can be useful to determine the areas in your network that may be susceptible to certain security risks or performance issues.

```sql+postgres
select
  name,
  zone_id,
  ttl
from
  cloudflare_load_balancer
where
  proxied;
```

```sql+sqlite
select
  name,
  zone_id,
  ttl
from
  cloudflare_load_balancer
where
  proxied = 1;
```

### Get session_affinity details for load balancers
Analyze the settings to understand the session affinity details of your load balancers. This can provide insights into how your web traffic is distributed and managed across multiple servers, aiding in efficient load balancing.

```sql+postgres
select
  name,
  zone_name,
  session_affinity,
  session_affinity_ttl,
  jsonb_pretty(session_affinity_attributes) as session_affinity_attributes
from
  cloudflare_load_balancer;
```

```sql+sqlite
select
  name,
  zone_name,
  session_affinity,
  session_affinity_ttl,
  session_affinity_attributes
from
  cloudflare_load_balancer;
```