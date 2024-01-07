package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
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
		var pathPattern, origin, viewerProtocolPolicy string
		cachePolicyId, originRequestPolicyId, responseHeadersPolicyId := "-", "-", "-"
		if v.PathPattern != nil {
			pathPattern = *v.PathPattern
		}
		if v.TargetOriginId != nil {
			origin = *v.TargetOriginId
		}
		viewerProtocolPolicy = utils.AutoCase(string(v.ViewerProtocolPolicy))
		if v.CachePolicyId != nil {
			cachePolicyId = *v.CachePolicyId
		}
		if v.OriginRequestPolicyId != nil {
			originRequestPolicyId = *v.OriginRequestPolicyId
		}
		if v.ResponseHeadersPolicyId != nil {
			responseHeadersPolicyId = *v.ResponseHeadersPolicyId
		}
		data = append(data, []string{
			pathPattern,
			origin,
			viewerProtocolPolicy,
			cachePolicyId,
			originRequestPolicyId,
			responseHeadersPolicyId,
		})
	}
	c.SetData(data)
}
