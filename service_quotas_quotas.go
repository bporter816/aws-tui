package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sq "github.com/aws/aws-sdk-go-v2/service/servicequotas"
	sqTypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type ServiceQuotasQuotas struct {
	*ui.Table
	sqClient    *sq.Client
	serviceName string
	serviceCode string
	app         *Application
}

func NewServiceQuotasQuotas(sqClient *sq.Client, serviceName string, serviceCode string, app *Application) *ServiceQuotasQuotas {
	s := &ServiceQuotasQuotas{
		Table: ui.NewTable([]string{
			"NAME",
			"APPLED VALUE",
			"DEFAULT VALUE",
			"ADJUSTABLE",
		}, 1, 0),
		sqClient:    sqClient,
		serviceName: serviceName,
		serviceCode: serviceCode,
		app:         app,
	}
	return s
}

func (s ServiceQuotasQuotas) GetService() string {
	return "Service Quotas"
}

func (s ServiceQuotasQuotas) GetLabels() []string {
	return []string{s.serviceName, "Quotas"}
}

func (s ServiceQuotasQuotas) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s ServiceQuotasQuotas) Render() {
	defaultsPg := sq.NewListAWSDefaultServiceQuotasPaginator(
		s.sqClient,
		&sq.ListAWSDefaultServiceQuotasInput{
			ServiceCode: aws.String(s.serviceCode),
		},
	)
	var defaultQuotas []sqTypes.ServiceQuota
	for defaultsPg.HasMorePages() {
		out, err := defaultsPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		defaultQuotas = append(defaultQuotas, out.Quotas...)
	}

	appliedPg := sq.NewListServiceQuotasPaginator(
		s.sqClient,
		&sq.ListServiceQuotasInput{
			ServiceCode: aws.String(s.serviceCode),
		},
	)
	appliedQuotas := make(map[string]sqTypes.ServiceQuota)
	for appliedPg.HasMorePages() {
		out, err := appliedPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, q := range out.Quotas {
			if q.QuotaName != nil {
				appliedQuotas[*q.QuotaName] = q
			}
		}
	}

	var data [][]string
	for _, v := range defaultQuotas {
		name, appliedValue, defaultValue, adjustable := "", "Not available", "", "No"
		if v.QuotaName != nil {
			name = *v.QuotaName
		}
		defaultValue = formatValue(v)
		if av, ok := appliedQuotas[name]; ok {
			appliedValue = formatValue(av)
		}
		if v.Adjustable {
			adjustable = "Yes"
		}
		data = append(data, []string{
			name,
			appliedValue,
			defaultValue,
			adjustable,
		})
	}
	s.SetData(data)
}

func formatValue(quota sqTypes.ServiceQuota) string {
	if quota.Value == nil {
		return ""
	}

	value := utils.SimplifyFloat(*quota.Value)

	if quota.Unit != nil && *quota.Unit != "None" {
		return fmt.Sprintf("%v %v", value, utils.AbbreviateUnit(*quota.Unit))
	} else {
		return value
	}
}
