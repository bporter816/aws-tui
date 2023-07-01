package main

import (
	"context"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type CloudfrontDistributions struct {
	*Table
	cfClient *cf.Client
	app      *Application
}

func NewCloudfrontDistributions(cfClient *cf.Client, app *Application) *CloudfrontDistributions {
	c := &CloudfrontDistributions{
		Table: NewTable([]string{
			"ID",
			"DESCRIPTION",
			"STATUS",
			"DOMAIN",
		}, 1, 0),
		cfClient: cfClient,
		app:      app,
	}
	c.Render() // TODO fix
	return c
}

func (c CloudfrontDistributions) GetName() string {
	return "Cloudfront"
}

func (c CloudfrontDistributions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CloudfrontDistributions) Render() {
	distributionsPaginator := cf.NewListDistributionsPaginator(c.cfClient, &cf.ListDistributionsInput{})
	var distributions []cfTypes.DistributionSummary
	for distributionsPaginator.HasMorePages() {
		out, err := distributionsPaginator.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		distributions = append(distributions, out.DistributionList.Items...)
	}

	var data [][]string
	for _, v := range distributions {
		var comment string
		if v.Comment != nil {
			comment = *v.Comment
		}
		data = append(data, []string{
			*v.Id,
			comment,
			*v.Status,
			*v.DomainName,
		})
	}
	c.SetData(data)
}
