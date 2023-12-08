---
title: "Steampipe Table: cloudflare_load_balancer_pool - Query Cloudflare Load Balancer Pools using SQL"
description: "Allows users to query Cloudflare Load Balancer Pools, specifically providing details about each load balancer pool, including its ID, created date, description, and enabled status."
---

# Table: cloudflare_load_balancer_pool - Query Cloudflare Load Balancer Pools using SQL

Cloudflare Load Balancer Pools are a component of Cloudflare's Load Balancing service. They consist of one or more origin servers where Cloudflare can direct traffic. These pools provide a way to distribute network or application traffic across many servers to optimize resource usage, reduce latency, and increase capacity to handle large amounts of traffic.

## Table Usage Guide

The `cloudflare_load_balancer_pool` table provides insights into Load Balancer Pools within Cloudflare's Load Balancing service. As a network administrator, explore pool-specific details through this table, including pool ID, creation date, description, and enabled status. Utilize it to uncover information about pools, such as their geographical distribution, the health of the servers within them, and the balance of traffic across the servers.

## Examples

### Basic info
Explore the status and details of Cloudflare load balancer pools, such as their names, identifiers, and whether they are enabled or not. This information can be useful in understanding and managing the distribution of network traffic across multiple servers.

```sql+postgres
select
  name,
  id,
  enabled,
  description,
  created_on,
  jsonb_pretty(origins) as origins
from
  cloudflare_load_balancer_pool;
```

```sql+sqlite
select
  name,
  id,
  enabled,
  description,
  created_on,
  origins
from
  cloudflare_load_balancer_pool;
```

### List active pools
Analyze the settings to understand the active load balancer pools in your Cloudflare account. This can help you manage your resources effectively by identifying which pools are currently in use.

```sql+postgres
select
  name,
  id,
  monitor
from
  cloudflare_load_balancer_pool
where
  enabled;
```

```sql+sqlite
select
  name,
  id,
  monitor
from
  cloudflare_load_balancer_pool
where
  enabled = 1;
```