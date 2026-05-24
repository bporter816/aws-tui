package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewServiceQuotasQuotas(repo *repo.ServiceQuotas, serviceName string, serviceCode string, app *Application) *TableView[model.ServiceQuotasQuota] {
	return NewTableView(
		"Service Quotas",
		[]string{serviceName, "Quotas"},
		[]string{"NAME", "APPLIED VALUE", "DEFAULT VALUE", "ADJUSTABLE"},
		func(ctx context.Context) ([]model.ServiceQuotasQuota, error) {
			return repo.ListQuotas(ctx, serviceCode)
		},
		func(v model.ServiceQuotasQuota) []string {
			appliedValue := "Not available"
			if v.AppliedValue != nil {
				appliedValue = utils.FormatServiceQuotasValue(v.AppliedValue, v.Unit)
			}
			return []string{
				utils.DerefString(v.QuotaName, ""),
				appliedValue,
				utils.FormatServiceQuotasValue(v.DefaultValue, v.Unit),
				utils.BoolToString(v.Adjustable, "Yes", "No"),
			}
		},
		nil,
	)
}
