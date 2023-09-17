package main

import (
	"fmt"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"strconv"
	"strings"
)

type Route53Records struct {
	*ui.Table
	repo         *repo.Route53
	hostedZoneId string
	app          *Application
}

func NewRoute53Records(repo *repo.Route53, zoneId string, app *Application) *Route53Records {
	r := &Route53Records{
		Table: ui.NewTable([]string{
			"RECORD NAME",
			"TYPE",
			// "ROUTING POLICY",
			"ROUTING",
			// "DIFFERENTIATOR",
			"DIFF",
			"LABEL",
			"TTL",
			"ALIAS",
			"VALUE",
		}, 1, 1),
		repo:         repo,
		hostedZoneId: zoneId,
		app:          app,
	}
	return r
}

func (r Route53Records) GetService() string {
	return "Route 53"
}

func (r Route53Records) GetLabels() []string {
	return []string{r.hostedZoneId, "Records"}
}

func (r Route53Records) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r Route53Records) Render() {
	model, err := r.repo.ListRecords(r.hostedZoneId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
			data = append(data, []string{
				strings.TrimSuffix(*v.Name, "."),
				string(v.Type),
				routingPolicy,
				differentiator,
				label,
				strconv.FormatInt(*v.TTL, 10),
				"No",
				// joinRoute53ResourceRecords(v.ResourceRecords, ","),
				fmtResourceRecords(v.ResourceRecords),
			})
		} else {
			// is an alias
			data = append(data, []string{
				strings.TrimSuffix(*v.Name, "."),
				string(v.Type),
				routingPolicy,
				differentiator,
				label,
				"-",
				"Yes",
				*v.AliasTarget.DNSName,
			})
		}
	}
	r.SetData(data)
}

func fmtResourceRecords(items []r53Types.ResourceRecord) string {
	if len(items) == 0 {
		return ""
	} else if len(items) == 1 {
		return *items[0].Value
	} else {
		return fmt.Sprintf("%v + %v more", *items[0].Value, len(items)-1)
	}
}

func joinRoute53ResourceRecords(items []r53Types.ResourceRecord, sep string) string {
	if len(items) == 0 {
		return ""
	}
	var ret string
	for _, v := range items {
		ret += sep
		ret += *v.Value
	}
	return strings.TrimPrefix(ret, sep)
}
