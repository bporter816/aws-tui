package internal

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewRoute53Records(repo *repo.Route53, zoneId string, app *Application) *TableView[model.Route53Record] {
	return NewTableView(
		"Route 53",
		[]string{zoneId, "Records"},
		[]string{"RECORD NAME", "TYPE", "ROUTING", "DIFF", "LABEL", "TTL", "ALIAS", "VALUE"},
		func() ([]model.Route53Record, error) {
			return repo.ListRecords(zoneId)
		},
		func(v model.Route53Record) []string {
			routingPolicy := "Simple"
			differentiator := "-"
			label := "-"

			// TODO IP routing policy and health checks

			if string(v.Failover) != "" {
				routingPolicy = "Failover"
				differentiator = string(v.Failover)
			}

			if string(v.Region) != "" {
				routingPolicy = "Latency"
				differentiator = string(v.Region)
			}

			if v.GeoLocation != nil {
				routingPolicy = "Geolocation"
				// TODO verify logic here
				if v.GeoLocation.ContinentCode != nil {
					differentiator = *v.GeoLocation.ContinentCode
				} else if v.GeoLocation.CountryCode != nil {
					differentiator = *v.GeoLocation.CountryCode
				} else if v.GeoLocation.SubdivisionCode != nil {
					// TODO include country code
					differentiator = *v.GeoLocation.SubdivisionCode
				}
			}

			if v.MultiValueAnswer != nil && *v.MultiValueAnswer {
				routingPolicy = "MultiValue"
				// TODO is there anything for the differentiator?
			}

			if v.Weight != nil {
				routingPolicy = "Weighted"
				differentiator = strconv.FormatInt(*v.Weight, 10)
			}

			if routingPolicy != "Simple" && v.SetIdentifier != nil {
				label = *v.SetIdentifier
			}

			if v.AliasTarget == nil {
				// not an alias
				return []string{
					strings.TrimSuffix(*v.Name, "."),
					string(v.Type),
					routingPolicy,
					differentiator,
					label,
					strconv.FormatInt(*v.TTL, 10),
					"No",
					// TODO consider removing, also see utils/route53.go
					// utils.JoinRoute53ResourceRecords(v.ResourceRecords, ","),
					utils.FormatRoute53ResourceRecords(v.ResourceRecords),
				}
			}
			// is an alias
			return []string{
				strings.TrimSuffix(*v.Name, "."),
				string(v.Type),
				routingPolicy,
				differentiator,
				label,
				"-",
				"Yes",
				*v.AliasTarget.DNSName,
			}
		},
		nil,
	)
}
