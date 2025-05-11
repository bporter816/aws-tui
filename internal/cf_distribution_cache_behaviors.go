package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type CFDistributionCacheBehaviors struct {
	*ui.Table
	view.CloudFront
	repo           *repo.CloudFront
	distributionId string
	app            *Application
}

func NewCFDistributionCacheBehaviors(repo *repo.CloudFront, distributionId string, app *Application) *CFDistributionCacheBehaviors {
	c := &CFDistributionCacheBehaviors{
		Table: ui.NewTable([]string{
			"PATH",
			"ORIGIN",
			"VIEWER PROTOCOL POLICY",
			"CACHE POLICY",
			"ORIGIN REQUEST POLICY",
			"RESPONSE HEADERS POLICY",
		}, 1, 0),
		repo:           repo,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionCacheBehaviors) GetLabels() []string {
	return []string{c.distributionId, "Cache Behaviors"}
}

func (c CFDistributionCacheBehaviors) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionCacheBehaviors) Render() {
	model, err := c.repo.GetDistributionCacheBehaviors(c.distributionId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.PathPattern, ""),
			utils.DerefString(v.TargetOriginId, ""),
			utils.AutoCase(string(v.ViewerProtocolPolicy)),
			utils.DerefString(v.CachePolicyId, "-"),
			utils.DerefString(v.OriginRequestPolicyId, "-"),
			utils.DerefString(v.ResponseHeadersPolicyId, "-"),
		})
	}
	c.SetData(data)
}
