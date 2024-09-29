// Copyright 2024 Jelly Terra
// Use of this source code form is governed under the MIT license.

package config

type AddrFilterRule struct {
	Ignore []string `yaml:"ignore"`
}

type NetIfAddrFilter struct {
	Ip4 *AddrFilterRule `yaml:"ip4"`
	Ip6 *AddrFilterRule `yaml:"ip6"`
}

type NetIf struct {
	Addr NetIfAddrFilter `yaml:"addr"`
}

type Record struct {
	Name    string `yaml:"name"`
	TTL     int    `yaml:"ttl"`
	Proxied bool   `yaml:"proxied"`

	Comment string `yaml:"comment"`

	NetIf map[string]*NetIf `yaml:"netif"`
}

type Zone struct {
	Token  string `yaml:"api_token"`
	ZoneID string `yaml:"zone_id"`

	MachineID string `yaml:"machine_id"`
	Shared    bool   `yaml:"shared"`

	Records []*Record `yaml:"records"`
}

type File struct {
	Zones map[string]*Zone `yaml:"zones"`
}
