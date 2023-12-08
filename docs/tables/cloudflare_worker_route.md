---
title: "Steampipe Table: cloudflare_worker_route - Query Cloudflare Worker Routes using SQL"
description: "Allows users to query Cloudflare Worker Routes, specifically the routes on which workers are triggered, providing insights into worker distribution and route coverage."
---

# Table: cloudflare_worker_route - Query Cloudflare Worker Routes using SQL

Cloudflare Worker Routes is a feature within Cloudflare that allows you to specify paths on your website where Cloudflare Workers scripts should be triggered. It provides a flexible way to control the behavior of your website by allowing you to run serverless code at the edge of Cloudflare's network, close to your users. Cloudflare Worker Routes helps you improve the performance and security of your web applications by executing scripts for specific routes.

## Table Usage Guide

The `cloudflare_worker_route` table provides insights into Worker Routes within Cloudflare. As a DevOps engineer, explore route-specific details through this table, including the paths where scripts are triggered and the corresponding script ID. Utilize it to uncover information about routes, such as those with specific scripts, the distribution of scripts across different routes, and the optimization of script execution.

## Examples

### Basic info
Explore which Cloudflare worker routes are currently in place. This can help identify potential areas for optimization or troubleshooting.

```sql+postgres
select
  id,
  zone_name,
  pattern,
  script
from
  cloudflare_worker_route;
```

```sql+sqlite
select
  id,
  zone_name,
  pattern,
  script
from
  cloudflare_worker_route;
```

### List idle worker routes (i.e. not attached to any worker)
Discover the segments that consist of idle worker routes that are not attached to any worker. This is beneficial for identifying unused resources and optimizing your Cloudflare configuration.

```sql+postgres
select
  id,
  zone_name,
  pattern
from
  cloudflare.cloudflare_worker_route
where
  script = '';
```

```sql+sqlite
select
  id,
  zone_name,
  pattern
from
  cloudflare_worker_route
where
  script = '';
```