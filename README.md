# Cloudflare Dynamic DNS
Dynamic DNS records updater implemented via Cloudflare API.

## Features

- Manage multiple zones.
- Update specified DNS records on event received from `netlink`.
- Address filter based on glob.

## Install

### Binary

Download from [releases](https://github.com/jellyterra/cloudflare-ddns/releases).

### Install via Go mod

```shell
$ go install github.com/jellyterra/cloudflare-ddns/cmd/cloudflare-ddns@latest
```

## Configure

### Zone

`[zones.(zoneName)]`

| Key          | Value  | Description                                        |
|--------------|--------|----------------------------------------------------|
| `zone_id`    | String | Zone ID of the Cloudflare Site.                    |
| `token`      | String | Cloudflare API token with accesses to the zone.    |
| `machine_id` | String | Machine ID, will be shown as comment on dashboard. |

### Record
`[zones.(zoneName).records.(domainName)]`

| Key       | Value  | Description         |
|-----------|--------|---------------------|
| `name`    | String | Domain name.        |
| `ttl`     | Int    | Record TTL seconds. |
| `proxied` | Bool   | Use Cloudflare CDN. |

### Network interface & Address filter

`[zones.(zoneName).records.(domainName).addr.(netifName).ip6]`

| Key      | Value    | Description                  |
|----------|----------|------------------------------|
| `ignore` | []String | Glob rules for IP addresses. |

You may have to filter reserved or internal IP sections via this feature.

Glob specification can be found [here](https://pkg.go.dev/regexp/syntax).

### Example

Example: `/etc/cloudflare-ddns.toml`
```toml
[zones.jellyterra]
zone_id = ""
token = ""
machine_id = "spacemit-k1"

[zones.jellyterra.records.k1]
name = "k1.jellyterra.com"
ttl = 60
proxied = true

[zones.jellyterra.records.k1.netif.wlp1s0.addr.ip4]
ignore = [ "192.*" ]

[zones.jellyterra.records.k1.netif.wlp1s0.addr.ip6]
ignore = [ "fe.*" ]
```

## Run

```shell
$ cloudflare-ddns -c /etc/cloudflare-ddns.toml
```
