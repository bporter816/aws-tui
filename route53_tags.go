package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53Tags struct {
	*Table
	r53Client    *r53.Client
	resourceType r53Types.TagResourceType
	resourceName string
	app          *Application
}

func NewRoute53Tags(r53Client *r53.Client, resourceType r53Types.TagResourceType, resourceName string, app *Application) *Route53Tags {
	r := &Route53Tags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		r53Client:    r53Client,
		resourceType: resourceType,
		resourceName: resourceName,
		app:          app,
	}
	return r
}

func (r Route53Tags) GetName() string {
	switch r.resourceType {
	case r53Types.TagResourceTypeHostedzone:
		return fmt.Sprintf("Route 53 | Hosted Zones | %v | Tags", r.resourceName)
	case r53Types.TagResourceTypeHealthcheck:
		return fmt.Sprintf("Route 53 | Health Checks | %v | Tags", r.resourceName)
	default:
		return fmt.Sprintf("Route 53 | <unknown> | %v | Tags", r.resourceName)
	}
}

func (r Route53Tags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r Route53Tags) Render() {
	out, err := r.r53Client.ListTagsForResource(
		context.TODO(),
		&r53.ListTagsForResourceInput{
			ResourceId:   aws.String(r.resourceName),
			ResourceType: r.resourceType,
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.ResourceTagSet != nil {
		for _, v := range out.ResourceTagSet.Tags {
			data = append(data, []string{
				*v.Key,
				*v.Value,
			})
		}
	}
	r.SetData(data)
}
