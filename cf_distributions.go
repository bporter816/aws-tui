package main

import (
	"context"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type CFDistributions struct {
	*Table
	cfClient *cf.Client
	app      *Application
}

func NewCFDistributions(cfClient *cf.Client, app *Application) *CFDistributions {
	c := &CFDistributions{
		Table: NewTable([]string{
			"ID",
			"DESCRIPTION",
			"STATUS",
			"DOMAIN",
		}, 1, 0),
		cfClient: cfClient,
		app:      app,
	}
	return c
}

func (c CFDistributions) GetName() string {
	return "Cloudfront | Distributions"
}

func (c CFDistributions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributions) Render() {
	pg := cf.NewListDistributionsPaginator(
		c.cfClient,
		&cf.ListDistributionsInput{},
	)
	var distributions []cfTypes.DistributionSummary
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
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
