package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type CFDistributionInvalidationPaths struct {
	*ui.Table
	view.CloudFront
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
