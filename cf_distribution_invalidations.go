package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type CFDistributionInvalidations struct {
	*ui.Table
	repo           *repo.CloudFront
	distributionId string
	app            *Application
}

func NewCFDistributionInvalidations(repo *repo.CloudFront, distributionId string, app *Application) *CFDistributionInvalidations {
	c := &CFDistributionInvalidations{
		Table: ui.NewTable([]string{
			"ID",
			"STATUS",
			"CREATED",
		}, 1, 0),
		repo:           repo,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionInvalidations) GetService() string {
	return "CloudFront"
}

func (c CFDistributionInvalidations) GetLabels() []string {
	return []string{c.distributionId, "Invalidations"}
}

func (c CFDistributionInvalidations) pathsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	pathsView := NewCFDistributionInvalidationPaths(c.repo, c.distributionId, id, c.app)
	c.app.AddAndSwitch(pathsView)
}

func (c CFDistributionInvalidations) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Paths",
			Action:      c.pathsHandler,
		},
	}
}

func (c CFDistributionInvalidations) Render() {
	model, err := c.repo.ListInvalidations(c.distributionId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
