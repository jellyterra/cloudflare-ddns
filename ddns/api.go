package ddns

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"net/netip"
)

func DeleteAllRecords(ctx context.Context, zone *Zone, recordName string) error {
	rc := &cloudflare.ResourceContainer{Identifier: zone.Raw.ZoneID}

	records, _, err := zone.API.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Name == recordName {

			if zone.Raw.Shared && record.Comment != zone.Raw.MachineID { // Owned by another host.
				continue
			}

			err = zone.API.DeleteDNSRecord(ctx, rc, record.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateRecord(ctx context.Context, zone *Zone, record *Record, typ, content string) error {
	_, err := zone.API.CreateDNSRecord(ctx, &cloudflare.ResourceContainer{Identifier: zone.Raw.ZoneID}, cloudflare.CreateDNSRecordParams{
		Type:    typ,
		Name:    record.Raw.Name,
		Content: content,
		TTL:     record.Raw.TTL,
		Proxied: &record.Raw.Proxied,
		Comment: zone.Raw.MachineID,
	})

	return err
}

func CreateIPv4Record(ctx context.Context, zone *Zone, record *Record, content string) error {
	return CreateRecord(ctx, zone, record, "A", content)
}

func CreateIPv6Record(ctx context.Context, zone *Zone, record *Record, content string) error {
	return CreateRecord(ctx, zone, record, "AAAA", content)
}

func UpdateRecord(ctx context.Context, zone *Zone, record *Record, addrs []netip.Addr) error {

	err := DeleteAllRecords(ctx, zone, record.Raw.Name)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		switch {
		case addr.Is4():
			err = CreateIPv4Record(ctx, zone, record, addr.String())
			if err != nil {
				return err
			}
		case addr.Is6():
			err = CreateIPv6Record(ctx, zone, record, addr.String())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
