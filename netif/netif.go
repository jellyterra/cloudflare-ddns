package netif

import (
	"net"
	"net/netip"
	"strings"
)

func GetNetIfAddrs(name string) (addrs []netip.Addr, _ error) {

	sysNetIf, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	rawAddrs, err := sysNetIf.Addrs()
	if err != nil {
		return nil, err
	}

	for _, rawAddr := range rawAddrs {
		addr, err := netip.ParseAddr(strings.Split(rawAddr.String(), "/")[0]) // Remove mask
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, addr)
	}

	return
}
