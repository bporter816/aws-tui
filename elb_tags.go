package main

import (
	"context"
	"fmt"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type ELBTags struct {
	*Table
	elbClient    *elb.Client
	resourceType ELBResourceType
	id           string
	name         string
	app          *Application
}

type ELBResourceType string

const (
	ELBResourceTypeLoadBalancer = "Load Balancers"
	ELBResourceTypeTargetGroup  = "Target Groups"
)

func NewELBTags(elbClient *elb.Client, resourceType ELBResourceType, id string, name string, app *Application) *ELBTags {
	e := &ELBTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		elbClient:    elbClient,
		resourceType: resourceType,
		id:           id,
		name:         name,
		app:          app,
	}
	return e
}

func (e ELBTags) GetName() string {
	return fmt.Sprintf("ELB | %v | %v | Tags", e.resourceType, e.name)
}

func (e ELBTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTags) Render() {
	out, err := e.elbClient.DescribeTags(
		context.TODO(),
		&elb.DescribeTagsInput{
			ResourceArns: []string{e.id},
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if len(out.TagDescriptions) != 1 {
		panic("should get exactly 1 tag description")
	}
	for _, v := range out.TagDescriptions[0].Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
