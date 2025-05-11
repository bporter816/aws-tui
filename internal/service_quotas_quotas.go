package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type ServiceQuotasQuotas struct {
	*ui.Table
	view.ServiceQuotas
	repo        *repo.ServiceQuotas
	serviceName string
	serviceCode string
	app         *Application
}

func NewServiceQuotasQuotas(repo *repo.ServiceQuotas, serviceName string, serviceCode string, app *Application) *ServiceQuotasQuotas {
	s := &ServiceQuotasQuotas{
		Table: ui.NewTable([]string{
			"NAME",
			"APPLIED VALUE",
			"DEFAULT VALUE",
			"ADJUSTABLE",
		}, 1, 0),
		repo:        repo,
		serviceName: serviceName,
		serviceCode: serviceCode,
		app:         app,
	}
	return s
}

func (s ServiceQuotasQuotas) GetLabels() []string {
	return []string{s.serviceName, "Quotas"}
}

func (s ServiceQuotasQuotas) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s ServiceQuotasQuotas) Render() {
	model, err := s.repo.ListQuotas(s.serviceCode)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		appliedValue := "Not available"
		if v.AppliedValue != nil {
			appliedValue = utils.FormatServiceQuotasValue(v.AppliedValue, v.Unit)
		}
		data = append(data, []string{
			utils.DerefString(v.QuotaName, ""),
			appliedValue,
			utils.FormatServiceQuotasValue(v.DefaultValue, v.Unit),
			utils.BoolToString(v.Adjustable, "Yes", "No"),
		})
	}
	s.SetData(data)
}
