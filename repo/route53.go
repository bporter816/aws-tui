package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bporter816/aws-tui/model"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53 struct {
	r53Client *r53.Client
}

func NewRoute53(r53Client *r53.Client) *Route53 {
	return &Route53{
		r53Client: r53Client,
	}
}

func (r Route53) ListRecords(hostedZoneId string) ([]model.Route53Record, error) {
	// ListResourceRecordSets doesn't have a paginator :'(
	good := true
	var resourceRecordSets []model.Route53Record
	var nextRecordName *string = nil
	var nextRecordType r53Types.RRType
	var nextRecordIdentifier *string = nil
	for good {
		out, err := r.r53Client.ListResourceRecordSets(
			context.TODO(),
			&r53.ListResourceRecordSetsInput{
				HostedZoneId:          aws.String(hostedZoneId),
				StartRecordName:       nextRecordName,
				StartRecordType:       nextRecordType,
				StartRecordIdentifier: nextRecordIdentifier,
			},
		)
		if err != nil {
			return []model.Route53Record{}, err
		}
		for _, v := range out.ResourceRecordSets {
			resourceRecordSets = append(resourceRecordSets, model.Route53Record(v))
		}
		good = out.IsTruncated
		if out.IsTruncated {
			nextRecordName = out.NextRecordName
			nextRecordType = out.NextRecordType
			nextRecordIdentifier = out.NextRecordIdentifier
		}
	}
	return resourceRecordSets, nil
}

func (r Route53) ListTags(resourceName string, resourceType r53Types.TagResourceType) (model.Tags, error) {
	out, err := r.r53Client.ListTagsForResource(
		context.TODO(),
		&r53.ListTagsForResourceInput{
			ResourceId:   aws.String(resourceName),
			ResourceType: resourceType,
		},
	)
	if err != nil || out.ResourceTagSet == nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.ResourceTagSet.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
