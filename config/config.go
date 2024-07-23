package config

type AddrFilterRule struct {
	Ignore []string `toml:"ignore"`
}

type NetIfAddrFilter struct {
	Ip4 *AddrFilterRule `toml:"ip4"`
	Ip6 *AddrFilterRule `toml:"ip6"`
}

type NetIf struct {
	Addr NetIfAddrFilter `toml:"addr"`
}

type Record struct {
	Name    string `toml:"name"`
	TTL     int    `toml:"ttl"`
	Proxied bool   `toml:"proxied"`

	Comment string `toml:"comment"`

	NetIf map[string]*NetIf `toml:"netif"`
}

type Zone struct {
	Token  string `toml:"token"`
	ZoneID string `toml:"zone_id"`

	MachineID string `toml:"machine_id"`
	Shared    bool   `toml:"shared"`

	Records map[string]*Record `toml:"records"`
}

type File struct {
	Zones map[string]*Zone `toml:"zones"`
}
