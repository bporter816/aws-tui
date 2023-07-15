package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/bporter816/aws-tui/ui"
)

type CFDistributionInvalidationPaths struct {
	*ui.Table
	cfClient       *cf.Client
	distributionId string
	invalidationId string
	app            *Application
}

func NewCFDistributionInvalidationPaths(cfClient *cf.Client, distributionId string, invalidationId string, app *Application) *CFDistributionInvalidationPaths {
	c := &CFDistributionInvalidationPaths{
		Table: ui.NewTable([]string{
			"PATH",
		}, 1, 0),
		cfClient:       cfClient,
		distributionId: distributionId,
		invalidationId: invalidationId,
		app:            app,
	}
	return c
}

func (c CFDistributionInvalidationPaths) GetService() string {
	return "Cloudfront"
}

func (c CFDistributionInvalidationPaths) GetLabels() []string {
	return []string{c.invalidationId, "Paths"}
}

func (c CFDistributionInvalidationPaths) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionInvalidationPaths) Render() {
	out, err := c.cfClient.GetInvalidation(
		context.TODO(),
		&cf.GetInvalidationInput{
			DistributionId: aws.String(c.distributionId),
			Id:             aws.String(c.invalidationId),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.Invalidation != nil && out.Invalidation.InvalidationBatch != nil && out.Invalidation.InvalidationBatch.Paths != nil {
		for _, v := range out.Invalidation.InvalidationBatch.Paths.Items {
			data = append(data, []string{
				v,
			})
		}
	}
	c.SetData(data)
}
