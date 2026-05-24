package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewCFDistributionCacheBehaviors(repo *repo.CloudFront, distributionId string, app *Application) *TableView[model.CloudFrontDistributionCacheBehavior] {
	return NewTableView(
		"CloudFront",
		[]string{distributionId, "Cache Behaviors"},
		[]string{"PATH", "ORIGIN", "VIEWER PROTOCOL POLICY", "CACHE POLICY", "ORIGIN REQUEST POLICY", "RESPONSE HEADERS POLICY"},
		func() ([]model.CloudFrontDistributionCacheBehavior, error) {
			return repo.GetDistributionCacheBehaviors(distributionId)
		},
		func(v model.CloudFrontDistributionCacheBehavior) []string {
			return []string{
				utils.DerefString(v.PathPattern, ""),
				utils.DerefString(v.TargetOriginId, ""),
				utils.AutoCase(string(v.ViewerProtocolPolicy)),
				utils.DerefString(v.CachePolicyId, "-"),
				utils.DerefString(v.OriginRequestPolicyId, "-"),
				utils.DerefString(v.ResponseHeadersPolicyId, "-"),
			}
		},
		nil,
	)
}
