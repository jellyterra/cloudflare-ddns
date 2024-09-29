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
| `api_token`  | String | Cloudflare API token with accesses to the zone.    |
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

Example: `/etc/cloudflare-ddns.yaml`

```yaml
zones:
  jellyterra: # Alias for identification.
    zone_id: ""
    api_token: ""
    machine_id: "spacemit-k1" # Tagging the records. The program does not change record with mismatched tag.
    records:
      - name: "k1.jellyterra.com"
        ttl: 60
        proxied: true
        netif:
          wlp1s0: # Network interface name
            addr:
              ip4:
                ignore:
                  - "192.*"
              ip6:
                ignore:
                  - "fe.*"
          enp1s0:
            addr:
              ip4:
                ignore:
                  - "192.*"
              ip6:
                ignore:
                  - "fe.*"
  symboltics:
    zone_id: ""
    api_token: ""
    # ...
```

## Run

```shell
$ cloudflare-ddns -c /etc/cloudflare-ddns.yaml
```
