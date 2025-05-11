package internal

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type Route53Records struct {
	*ui.Table
	view.Route53
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
				// TODO consider removing, also see utils/route53.go
				// utils.JoinRoute53ResourceRecords(v.ResourceRecords, ","),
				utils.FormatRoute53ResourceRecords(v.ResourceRecords),
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
