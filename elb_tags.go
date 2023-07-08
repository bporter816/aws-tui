package main

import (
	"context"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type ELBTags struct {
	*Table
	elbClient    *elb.Client
	resourceType ELBResourceType
	resourceArn  string
	resourceName string
	app          *Application
}

type ELBResourceType string

const (
	ELBResourceTypeLoadBalancer ELBResourceType = "Load Balancers"
	ELBResourceTypeTargetGroup  ELBResourceType = "Target Groups"
)

func NewELBTags(elbClient *elb.Client, resourceType ELBResourceType, resourceArn string, resourceName string, app *Application) *ELBTags {
	e := &ELBTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		elbClient:    elbClient,
		resourceType: resourceType,
		resourceArn:  resourceArn,
		resourceName: resourceName,
		app:          app,
	}
	return e
}

func (e ELBTags) GetService() string {
	return "ELB"
}

func (e ELBTags) GetLabels() []string {
	return []string{e.resourceName, "Tags"}
}

func (e ELBTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTags) Render() {
	out, err := e.elbClient.DescribeTags(
		context.TODO(),
		&elb.DescribeTagsInput{
			ResourceArns: []string{e.resourceArn},
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
