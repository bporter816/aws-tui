package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
)

func NewCFDistributionInvalidationPaths(repo *repo.CloudFront, distributionId string, invalidationId string, app *Application) *TableView[model.CloudFrontInvalidationPath] {
	return NewTableView(
		"CloudFront",
		[]string{invalidationId, "Paths"},
		[]string{"PATH"},
		func() ([]model.CloudFrontInvalidationPath, error) {
			return repo.ListInvalidationPaths(distributionId, invalidationId)
		},
		func(v model.CloudFrontInvalidationPath) []string {
			return []string{string(v)}
		},
		nil,
	)
}
