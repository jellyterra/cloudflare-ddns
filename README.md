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

## Setup

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

```shell
$ cloudflare-ddns -c /etc/cloudflare-ddns.toml
```
