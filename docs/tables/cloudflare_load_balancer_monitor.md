---
title: "Steampipe Table: cloudflare_load_balancer_monitor - Query Cloudflare Load Balancer Monitors using SQL"
description: "Allows users to query Cloudflare Load Balancer Monitors, specifically their configuration and status, providing insights into the health and performance of the load balancers."
---

# Table: cloudflare_load_balancer_monitor - Query Cloudflare Load Balancer Monitors using SQL

Cloudflare Load Balancer Monitors are a feature of Cloudflare's Load Balancing service. They provide continuous checks on your servers to determine their health and direct traffic accordingly. Monitors help in managing traffic distribution, reducing latency, and improving data delivery speed.

## Table Usage Guide

The `cloudflare_load_balancer_monitor` table provides insights into the configuration and status of Load Balancer Monitors within Cloudflare. As a Network Administrator, explore monitor-specific details through this table, including type, method, path, and timeout settings. Utilize it to uncover information about monitors, such as their current status, frequency of health checks, and the expected codes for successful checks.

## Examples

### Basic info
Explore which Cloudflare load balancer monitors have certain configurations to optimize performance and reliability. This can help in identifying any monitors that may need adjustments for better load balancing and error handling.

```sql+postgres
select
  id,
  type,
  path,
  timeout,
  retries,
  interval,
  port,
  expected_codes
from
  cloudflare_load_balancer_monitor;
```

```sql+sqlite
select
  id,
  type,
  path,
  timeout,
  retries,
  interval,
  port,
  expected_codes
from
  cloudflare_load_balancer_monitor;
```

### Get information of monitors attached to pool
Explore which monitors are attached to specific load balancer pools in Cloudflare. This query is useful for gaining insights into the configuration and status of your load balancing setup, including details like pool ID, name, and notification email.

```sql+postgres
select
  p.id as pool_id,
  p.name as pool_name,
  p.enabled as pool_enabled,
  p.notification_email,
  m.id as monitor_id,
  m.description monitor_description
from
  cloudflare_load_balancer_pool p,
  cloudflare_load_balancer_monitor as m
where
  p.monitor = m.id;
```

```sql+sqlite
select
  p.id as pool_id,
  p.name as pool_name,
  p.enabled as pool_enabled,
  p.notification_email,
  m.id as monitor_id,
  m.description as monitor_description
from
  cloudflare_load_balancer_pool p,
  cloudflare_load_balancer_monitor m
where
  p.monitor = m.id;
```