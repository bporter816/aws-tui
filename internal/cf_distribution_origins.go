package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewCFDistributionOrigins(repo *repo.CloudFront, distributionId string, app *Application) *TableView[model.CloudFrontDistributionOrigin] {
	return NewTableView(
		"CloudFront",
		[]string{distributionId, "Origins"},
		[]string{"NAME", "DOMAIN NAME", "PATH", "TYPE"},
		func(ctx context.Context) ([]model.CloudFrontDistributionOrigin, error) {
			return repo.GetDistributionOrigins(ctx, distributionId)
		},
		func(v model.CloudFrontDistributionOrigin) []string {
			var originType string
			if v.S3OriginConfig != nil {
				originType = "S3"
			} else if v.CustomOriginConfig != nil {
				originType = "Custom origin"
			}
			return []string{
				utils.DerefString(v.Id, ""),
				utils.DerefString(v.DomainName, ""),
				utils.DerefString(v.OriginPath, ""),
				originType,
			}
		},
		nil,
	)
}
