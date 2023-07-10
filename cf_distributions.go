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

func (c CFDistributions) GetService() string {
	return "Cloudfront"
}

func (c CFDistributions) GetLabels() []string {
	return []string{"Distributions"}
}

func (c CFDistributions) originsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	originsView := NewCFDistributionOrigins(c.cfClient, id, c.app)
	c.app.AddAndSwitch(originsView)
}

func (c CFDistributions) cacheBehaviorsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	cacheBehaviorsView := NewCFDistributionCacheBehaviors(c.cfClient, id, c.app)
	c.app.AddAndSwitch(cacheBehaviorsView)
}

func (c CFDistributions) customErrorResponsesHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	cacheBehaviorsView := NewCFDistributionCustomErrorResponses(c.cfClient, id, c.app)
	c.app.AddAndSwitch(cacheBehaviorsView)
}

func (c CFDistributions) tagsHandler() {
	row, err := c.GetRowSelection()
	if err != nil {
		return
	}
	tagsView := NewCFTags(c.cfClient, c.arns[row-1], c.app)
	c.app.AddAndSwitch(tagsView)
}

func (c CFDistributions) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone),
			Description: "Origins",
			Action:      c.originsHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "Cache Behaviors",
			Action:      c.cacheBehaviorsHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
			Description: "Custom Error Responses",
			Action:      c.customErrorResponsesHandler,
		},
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
