---
organization: Turbot
category: ["internet"]
icon_url: "/images/plugins/turbot/cloudflare.svg"
brand_color: "#f48120"
display_name: Cloudflare
name: cloudflare
description: Steampipe plugin for querying Cloudflare databases, networks, and other resources.
og_description: Query cloudflare databases, networks, and other resources with SQL! Open source CLI. No DB required. 
og_image: "/images/plugins/turbot/cloudflare-social-graphic.png"
---

# Cloudflare

Query your Cloudflare infrastructure including zones, DNS records, accounts and more.

## Installation

Download and install the latest Cloudflare plugin:

```bash
steampipe plugin install cloudflare
```

## Configuration

Connection configurations are defined using HCL in one or more Steampipe config files. Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.

Installing the latest cloudflare plugin will create a connection file (`~/.steampipe/config/cloudflare.spc`) with a single connection named `cloudflare`. You must modify this connection to include your personal credentials.

An [API Token](https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys#12345680) is the recommended way to set credentials. Read scope is required (write is not):

```hcl
connection "cloudflare" {
  plugin  = "cloudflare"
  token   = "psth3GX0qHavRYE-hd5y7_iL7piII6C8jR3FOuW3"
}
```

It's also valid to use an email and API key:

```hcl
connection "cloudflare" {
  plugin  = "cloudflare"
  email   = "pam@dundermifflin.com"
  api_key = "2980b99351d629a537f1440e12b5b97a135b7"
}
```

Credentials are resolved in this order:

1. `token` in Steampipe config.
2. `email` and `api_key` in Steampipe config.
3. `CLOUDFLARE_API_TOKEN` environment variable (like Terraform).
4. `CLOUDFLARE_EMAIL` and `CLOUDFLARE_API_KEY` environment variables (like Terraform).
5. `CF_API_TOKEN` environment variable (like flarectl).
6. `CF_API_EMAIL` and `CF_API_KEY` environment variables (like flarectl).

For example:

```hcl
connection "cloudflare" {
  plugin = "cloudflare"
  token  = "9wZVRX3j9Z1CiE38HcmThwkb2hThisIsAFakeToken"
}
```

## Scope

A Cloudflare connection is scoped to a single Cloudflare account, with a single set of credentials.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-cloudflare
- Community: [Slack Channel](https://steampipe.io/community/join)
