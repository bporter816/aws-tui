package main

import (
	"context"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/gdamore/tcell/v2"
)

type CFDistributions struct {
	*Table
	cfClient *cf.Client
	app      *Application
	arns     []string
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

func (c CFDistributions) tagsHandler() {
	row, _ := c.GetSelection()
	tagsView := NewCFTags(c.cfClient, c.arns[row-1], c.app)
	c.app.AddAndSwitch(tagsView)
}

func (c CFDistributions) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      c.tagsHandler,
		},
	}
}

func (c *CFDistributions) Render() {
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
	c.arns = make([]string, len(distributions))
	for i, v := range distributions {
		c.arns[i] = *v.ARN
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
