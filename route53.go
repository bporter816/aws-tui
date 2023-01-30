package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"strconv"
)

type Route53 struct {
	client *route53.Client
}

func NewRoute53(c *route53.Client) *Route53 {
	return &Route53{
		client: c,
	}
}

func (r Route53) GetHeaders() []string {
	return []string{
		"ID",
		"NAME",
		"RECORDS",
		"DESCRIPTION",
		"VISIBILITY",
	}
}

func (r *Route53) Render() ([][]string, error) {
	hostedZonesPaginator := route53.NewListHostedZonesPaginator(r.client, &route53.ListHostedZonesInput{})
	var hostedZones []route53Types.HostedZone
	for hostedZonesPaginator.HasMorePages() {
		out, err := hostedZonesPaginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		hostedZones = append(hostedZones, out.HostedZones...)
	}

	var ret [][]string
	for _, v := range hostedZones {
		var comment, visibility string
		if v.Config != nil {
			if v.Config.Comment != nil {
				comment = *v.Config.Comment
			}
			if v.Config.PrivateZone {
				visibility = "Private"
			} else {
				visibility = "Public"
			}
		}
		ret = append(ret, []string{
			*v.Id,
			*v.Name,
			strconv.FormatInt(*v.ResourceRecordSetCount, 10),
			comment,
			visibility,
		})
	}
	return ret, nil
}
