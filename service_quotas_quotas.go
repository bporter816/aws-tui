package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sq "github.com/aws/aws-sdk-go-v2/service/servicequotas"
	sqTypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/bporter816/aws-tui/ui"
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
	appliedQuotas := make(map[string]float64)
	for appliedPg.HasMorePages() {
		out, err := appliedPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, q := range out.Quotas {
			if q.QuotaName != nil && q.Value != nil {
				appliedQuotas[*q.QuotaName] = *q.Value
			}
		}
	}

	var data [][]string
	for _, v := range defaultQuotas {
		name, appliedValue, defaultValue, adjustable := "", "Not available", "", "No"
		if v.QuotaName != nil {
			name = *v.QuotaName
		}
		var unit string
		if v.Value != nil {
			if v.Unit != nil && *v.Unit != "None" {
				unit = " " + *v.Unit
			}
			defaultValue = fmt.Sprintf("%f%v", *v.Value, unit)
		}
		if av, ok := appliedQuotas[name]; ok {
			appliedValue = fmt.Sprintf("%f%v", av, unit)
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
