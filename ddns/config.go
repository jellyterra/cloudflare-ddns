// Copyright 2024 Jelly Terra
// Use of this source code form is governed under the MIT license.

package ddns

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/jellyterra/cloudflare-ddns/config"
	"regexp"
)

func CompileRegExps(rawRegExps []string) ([]*regexp.Regexp, error) {
	regExps := make([]*regexp.Regexp, len(rawRegExps))
	for i, rawRegExp := range rawRegExps {
		regExp, err := regexp.Compile(rawRegExp)
		if err != nil {
			return nil, err
		}
		regExps[i] = regExp
	}
	return regExps, nil
}

func LoadNetIfAddrFilter(netIfName string, netIf *config.NetIf) (*NetIfAddrFilter, error) {
	var ip4Ignore, ip6Ignore []*regexp.Regexp

	if netIf.Addr.Ip4 != nil {
		regExps, err := CompileRegExps(netIf.Addr.Ip4.Ignore)
		if err != nil {
			return nil, err
		}
		ip4Ignore = regExps
	}

	if netIf.Addr.Ip6 != nil {
		regExps, err := CompileRegExps(netIf.Addr.Ip6.Ignore)
		if err != nil {
			return nil, err
		}
		ip6Ignore = regExps
	}

	return &NetIfAddrFilter{
		NetIfName: netIfName,
		Ip4Ignore: ip4Ignore,
		Ip6Ignore: ip6Ignore,
	}, nil
}

func LoadRecord(record *config.Record) (*Record, error) {
	var filters []*NetIfAddrFilter

	for netIfName, rawFilter := range record.NetIf {
		filter, err := LoadNetIfAddrFilter(netIfName, rawFilter)
		if err != nil {
			return nil, err
		}

		filters = append(filters, filter)
	}

	return &Record{
		Raw:              record,
		NetIfAddrFilters: filters,
	}, nil
}

func LoadZone(zoneKey string, zone *config.Zone) (*Zone, error) {

	var (
		records []*Record
	)

	for _, rawRecord := range zone.Records {
		record, err := LoadRecord(rawRecord)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	api, err := cloudflare.NewWithAPIToken(zone.Token)
	if err != nil {
		return nil, err
	}

	return &Zone{
		Raw:     zone,
		API:     api,
		Records: records,
		ZoneKey: zoneKey,
	}, nil
}

func LoadConfig(file *config.File) (*Env, error) {
	var zones []*Zone

	for zoneKey, rawZone := range file.Zones {
		zone, err := LoadZone(zoneKey, rawZone)
		if err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	return &Env{
		Zones: zones,
	}, nil
}
