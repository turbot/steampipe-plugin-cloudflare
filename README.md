![image](https://hub.steampipe.io/images/plugins/turbot/cloudflare-social-graphic.png)

# Cloudflare Plugin for Steampipe

Use SQL to query accounts, zones and more from Cloudflare.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/cloudflare)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/cloudflare/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-cloudflare/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install cloudflare
```

Run a query:

```sql
select
  name,
  dnssec ->> 'status',
  settings ->> 'tls_1_3'
from
  cloudflare_zone
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-cloudflare.git
cd steampipe-plugin-cloudflare
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/cloudflare.spc
```

Try it!

```
steampipe query
> .inspect cloudflare
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). Contributions to the plugin are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-cloudflare/blob/main/LICENSE). Contributions to the plugin documentation are subject to the [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-cloudflare/blob/main/docs/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Cloudflare Plugin](https://github.com/turbot/steampipe-plugin-cloudflare/labels/help%20wanted)
