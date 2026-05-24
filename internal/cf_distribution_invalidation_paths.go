package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
)

func NewCFDistributionInvalidationPaths(repo *repo.CloudFront, distributionId string, invalidationId string, app *Application) *TableView[model.CloudFrontInvalidationPath] {
	return NewTableView(
		"CloudFront",
		[]string{invalidationId, "Paths"},
		[]string{"PATH"},
		func(ctx context.Context) ([]model.CloudFrontInvalidationPath, error) {
			return repo.ListInvalidationPaths(ctx, distributionId, invalidationId)
		},
		func(v model.CloudFrontInvalidationPath) []string {
			return []string{string(v)}
		},
		nil,
	)
}
