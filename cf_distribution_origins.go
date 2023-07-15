package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/bporter816/aws-tui/ui"
)

type CFDistributionOrigins struct {
	*ui.Table
	cfClient       *cf.Client
	distributionId string
	app            *Application
}

func NewCFDistributionOrigins(cfClient *cf.Client, distributionId string, app *Application) *CFDistributionOrigins {
	c := &CFDistributionOrigins{
		Table: ui.NewTable([]string{
			"NAME",
			"DOMAIN NAME",
			"PATH",
			"TYPE",
		}, 1, 0),
		cfClient:       cfClient,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionOrigins) GetService() string {
	return "Cloudfront"
}

func (c CFDistributionOrigins) GetLabels() []string {
	return []string{c.distributionId, "Origins"}
}

func (c CFDistributionOrigins) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionOrigins) Render() {
	out, err := c.cfClient.GetDistributionConfig(
		context.TODO(),
		&cf.GetDistributionConfigInput{
			Id: aws.String(c.distributionId),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.DistributionConfig != nil && out.DistributionConfig.Origins != nil {
		for _, v := range out.DistributionConfig.Origins.Items {
			var id, domainName, originPath, originType string
			if v.Id != nil {
				id = *v.Id
			}
			if v.DomainName != nil {
				domainName = *v.DomainName
			}
			if v.OriginPath != nil {
				originPath = *v.OriginPath
			}
			if v.S3OriginConfig != nil {
				originType = "S3"
			} else if v.CustomOriginConfig != nil {
				originType = "Custom origin"
			}
			data = append(data, []string{
				id,
				domainName,
				originPath,
				originType,
			})
		}
	}
	c.SetData(data)
}
