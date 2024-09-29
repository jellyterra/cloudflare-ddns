// Copyright 2024 Jelly Terra
// Use of this source code form is governed under the MIT license.

package ddns

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/jellyterra/cloudflare-ddns/config"
	"github.com/jellyterra/cloudflare-ddns/netif"
	"log"
	"net/netip"
	"regexp"
	"sync"
)

func FilterAddrs(addrs []netip.Addr, patterns []*regexp.Regexp) (filtered []netip.Addr) {
	for _, addr := range addrs {
		matched := false
		for _, pattern := range patterns {
			if pattern.MatchString(addr.String()) {
				matched = true
				break
			}
		}
		if !matched {
			filtered = append(filtered, addr)
		}
	}
	return
}

type NetIfAddrFilter struct {
	NetIfName string

	Ip4Ignore, Ip6Ignore []*regexp.Regexp
}

func (f *NetIfAddrFilter) Filter() ([]netip.Addr, error) {
	var (
		ip4Addrs, ip6Addrs []netip.Addr
	)

	addrs, err := netif.GetNetIfAddrs(f.NetIfName)
	if err != nil {
		log.Println("Skipped interface", f.NetIfName, "due to:", err)
		return nil, nil // Skip.
	}

	for _, addr := range addrs {
		switch {
		case addr.Is4():
			ip4Addrs = append(ip4Addrs, addr)
		case addr.Is6():
			ip6Addrs = append(ip6Addrs, addr)
		}
	}

	ip4Addrs = FilterAddrs(ip4Addrs, f.Ip4Ignore)
	ip6Addrs = FilterAddrs(ip6Addrs, f.Ip6Ignore)

	return append(ip4Addrs, ip6Addrs...), nil
}

type Record struct {
	Raw *config.Record

	NetIfAddrFilters []*NetIfAddrFilter
}

type Zone struct {
	Raw *config.Zone

	API *cloudflare.API

	Records []*Record

	ZoneKey string
}

func (z *Zone) UpdateRecord(ctx context.Context, record *Record) error {

	var allAddrs []netip.Addr

	for _, filter := range record.NetIfAddrFilters {
		addrs, err := filter.Filter()
		if err != nil {
			return err
		}

		allAddrs = append(allAddrs, addrs...)
	}

	return UpdateRecord(ctx, z, record, allAddrs)
}

func (z *Zone) UpdateAll(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, record := range z.Records {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := z.UpdateRecord(ctx, record)
			if err != nil {
				log.Println("Record", record.Raw.Name, "update failed:", err)
			}
		}()
	}

	wg.Wait()

	return nil
}

type Env struct {
	Zones []*Zone
}

func (e *Env) UpdateAllZones(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, zone := range e.Zones {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := zone.UpdateAll(ctx)
			if err != nil {
				log.Println("Zone", zone.ZoneKey, "update failed:", err)
			}
		}()
	}

	wg.Wait()

	return nil
}
