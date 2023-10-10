package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type CFDistributionInvalidationPaths struct {
	*ui.Table
	repo           *repo.CloudFront
	distributionId string
	invalidationId string
	app            *Application
}

func NewCFDistributionInvalidationPaths(repo *repo.CloudFront, distributionId string, invalidationId string, app *Application) *CFDistributionInvalidationPaths {
	c := &CFDistributionInvalidationPaths{
		Table: ui.NewTable([]string{
			"PATH",
		}, 1, 0),
		repo:           repo,
		distributionId: distributionId,
		invalidationId: invalidationId,
		app:            app,
	}
	return c
}

func (c CFDistributionInvalidationPaths) GetService() string {
	return "CloudFront"
}

func (c CFDistributionInvalidationPaths) GetLabels() []string {
	return []string{c.invalidationId, "Paths"}
}

func (c CFDistributionInvalidationPaths) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionInvalidationPaths) Render() {
	model, err := c.repo.ListInvalidationPaths(c.distributionId, c.invalidationId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			string(v),
		})
	}
	c.SetData(data)
}
