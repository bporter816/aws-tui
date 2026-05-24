package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewELBTrustStoreAssociations(repo *repo.ELB, a arn.ARN, app *Application) *TableView[model.ELBTrustStoreAssociation] {
	return NewTableView(
		"ELB",
		[]string{a.Resource, "Associations"},
		[]string{"ARN"},
		func(ctx context.Context) ([]model.ELBTrustStoreAssociation, error) {
			return repo.ListTrustStoreAssociations(ctx, a)
		},
		func(v model.ELBTrustStoreAssociation) []string {
			return []string{
				utils.DerefString(v.ResourceArn, ""),
			}
		},
		nil,
	)
}
