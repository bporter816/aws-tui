package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"strconv"
	"strings"
)

type Route53Records struct {
	*Table
	r53Client    *route53.Client
	hostedZoneId string
}

func NewRoute53Records(client *route53.Client, zoneId string) *Route53Records {
	r := &Route53Records{
		Table: NewTable([]string{
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
		r53Client:    client,
		hostedZoneId: zoneId,
	}
	return r
}

func (r Route53Records) GetName() string {
	return "Route 53 Records"
}

func (r Route53Records) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r Route53Records) Render() {
	// ListResourceRecordSets doesn't have a paginator :'(
	good := true
	var resourceRecordSets []route53Types.ResourceRecordSet
	var nextRecordName *string = nil
	var nextRecordType route53Types.RRType
	var nextRecordIdentifier *string = nil
	for good {
		out, err := r.r53Client.ListResourceRecordSets(
			context.TODO(),
			&route53.ListResourceRecordSetsInput{
				HostedZoneId:          aws.String(r.hostedZoneId),
				StartRecordName:       nextRecordName,
				StartRecordType:       nextRecordType,
				StartRecordIdentifier: nextRecordIdentifier,
			},
		)
		if err != nil {
			panic(err)
		}
		resourceRecordSets = append(resourceRecordSets, out.ResourceRecordSets...)
		good = out.IsTruncated
		if out.IsTruncated {
			nextRecordName = out.NextRecordName
			nextRecordType = out.NextRecordType
			nextRecordIdentifier = out.NextRecordIdentifier
		}
	}

	var data [][]string
	for _, v := range resourceRecordSets {
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

		if routingPolicy != "Simple" {
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

func fmtResourceRecords(items []route53Types.ResourceRecord) string {
	if len(items) == 0 {
		return ""
	} else if len(items) == 1 {
		return *items[0].Value
	} else {
		return fmt.Sprintf("%v + %v more", *items[0].Value, len(items)-1)
	}
}

func joinRoute53ResourceRecords(items []route53Types.ResourceRecord, sep string) string {
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
