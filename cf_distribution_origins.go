package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

type CFDistributionOrigins struct {
	*Table
	cfClient       *cf.Client
	distributionId string
	app            *Application
}

func NewCFDistributionOrigins(cfClient *cf.Client, distributionId string, app *Application) *CFDistributionOrigins {
	c := &CFDistributionOrigins{
		Table: NewTable([]string{
			"NAME",
		}, 1, 0),
		cfClient:       cfClient,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionOrigins) GetName() string {
	return fmt.Sprintf("Cloudfront | Distributions | %v | Origins", c.distributionId)
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
			data = append(data, []string{
				*v.Id,
			})
		}
	}
	c.SetData(data)
}
