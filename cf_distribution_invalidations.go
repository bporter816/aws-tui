package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type CFDistributionInvalidations struct {
	*ui.Table
	cfClient       *cf.Client
	distributionId string
	app            *Application
}

func NewCFDistributionInvalidations(cfClient *cf.Client, distributionId string, app *Application) *CFDistributionInvalidations {
	c := &CFDistributionInvalidations{
		Table: ui.NewTable([]string{
			"ID",
			"STATUS",
			"CREATED",
		}, 1, 0),
		cfClient:       cfClient,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionInvalidations) GetService() string {
	return "Cloudfront"
}

func (c CFDistributionInvalidations) GetLabels() []string {
	return []string{c.distributionId, "Invalidations"}
}

func (c CFDistributionInvalidations) pathsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	pathsView := NewCFDistributionInvalidationPaths(c.cfClient, c.distributionId, id, c.app)
	c.app.AddAndSwitch(pathsView)
}

func (c CFDistributionInvalidations) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Paths",
			Action:      c.pathsHandler,
		},
	}
}

func (c CFDistributionInvalidations) Render() {
	pg := cf.NewListInvalidationsPaginator(
		c.cfClient,
		&cf.ListInvalidationsInput{
			DistributionId: aws.String(c.distributionId),
		},
	)
	var invalidations []cfTypes.InvalidationSummary
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		if out.InvalidationList != nil {
			invalidations = append(invalidations, out.InvalidationList.Items...)
		}
	}

	var data [][]string
	for _, v := range invalidations {
		id, status, created := "-", "-", "-"
		if v.Id != nil {
			id = *v.Id
		}
		if v.Status != nil {
			status = *v.Status
		}
		if v.CreateTime != nil {
			created = v.CreateTime.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			id,
			status,
			created,
		})
	}
	c.SetData(data)
}
