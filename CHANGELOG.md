## v0.0.3 [2021-10-04]

_Enhancements_

- Updated the README file as per the latest format ([#18](https://github.com/turbot/steampipe-plugin-cloudflare/pull/18))

_Bug fixes_
- Fixed the `cloudflare_zone` table to include the `dnssec` and `settings` columns ([#17](https://github.com/turbot/steampipe-plugin-cloudflare/pull/17))

## v0.0.2 [2021-07-02]

_What's new?_

- New tables added
  - [cloudflare_account_role](https://hub.steampipe.io/plugins/turbot/cloudflare/tables/cloudflare_account_role) ([#6](https://github.com/turbot/steampipe-plugin-cloudflare/pull/6))
  - [cloudflare_firewall_rule](https://hub.steampipe.io/plugins/turbot/cloudflare/tables/cloudflare_firewall_rule) ([#8](https://github.com/turbot/steampipe-plugin-cloudflare/pull/8))
  - [cloudflare_page_rule](https://hub.steampipe.io/plugins/turbot/cloudflare/tables/cloudflare_page_rule) ([#10](https://github.com/turbot/steampipe-plugin-cloudflare/pull/10))

_Enhancements_

- Updated plugin category from `networking` to `internet`
- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)

## v0.0.1 [2021-03-11]

_What's new?_

- Initial release with tables for accounts, API tokens, DNS records, users, and zones.
