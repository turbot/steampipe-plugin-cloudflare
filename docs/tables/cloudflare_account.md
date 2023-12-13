---
title: "Steampipe Table: cloudflare_account - Query Cloudflare Accounts using SQL"
description: "Allows users to query Cloudflare Accounts, specifically the account details such as ID, name, email, and status."
---

# Table: cloudflare_account - Query Cloudflare Accounts using SQL

Cloudflare is a web infrastructure and website security company that provides content delivery network services, DDoS mitigation, Internet security, and distributed domain name server services. Cloudflare's services sit between a website's visitor and the Cloudflare user's hosting provider, acting as a reverse proxy for websites. Its network protects, speeds up, and improves availability for a website or mobile application with a change in DNS.

## Table Usage Guide

The `cloudflare_account` table provides insights into accounts within Cloudflare. As a DevOps engineer or a security analyst, explore account-specific details through this table, including account ID, name, email, and status. Utilize it to uncover information about accounts, such as those with specific status, the verification of email addresses, and the identification of account names.

## Examples

### Query all accounts the user has access to
Determine the range of accounts to which a user has access. This allows for a comprehensive overview of user permissions, aiding in account management and security audits.

```sql+postgres
select
  *
from
  cloudflare_account;
```

```sql+sqlite
select
  *
from
  cloudflare_account;
```

### Check if two factor authentication is enforced for accounts
Analyze the settings to understand whether two-factor authentication is being enforced for accounts, thereby enhancing security measures.

```sql+postgres
select
  name,
  settings -> 'enforce_twofactor' as enforce_mfa
from
  cloudflare_account;
```

```sql+sqlite
select
  name,
  json_extract(settings, '$.enforce_twofactor') as enforce_mfa
from
  cloudflare_account;
```