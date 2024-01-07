package main

import (
	"fmt"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
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
		name, appliedValue, defaultValue, adjustable := "", "Not available", "", ""
		if v.QuotaName != nil {
			name = *v.QuotaName
		}
		defaultValue = formatValue(v.DefaultValue, v.Unit)
		if v.AppliedValue != nil {
			appliedValue = formatValue(v.AppliedValue, v.Unit)
		}
		adjustable = utils.BoolToString(v.Adjustable, "Yes", "No")
		data = append(data, []string{
			name,
			appliedValue,
			defaultValue,
			adjustable,
		})
	}
	s.SetData(data)
}

func formatValue(value *float64, unit *string) string {
	if value == nil {
		return ""
	}

	v := utils.SimplifyFloat(*value)

	if unit != nil && *unit != "None" {
		return fmt.Sprintf("%v %v", v, utils.AbbreviateUnit(*unit))
	} else {
		return v
	}
}
