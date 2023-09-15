package main

import (
	"context"
	sq "github.com/aws/aws-sdk-go-v2/service/servicequotas"
	sqTypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type ServiceQuotasServices struct {
	*ui.Table
	sqClient *sq.Client
	app      *Application
}

func NewServiceQuotasServices(sqClient *sq.Client, app *Application) *ServiceQuotasServices {
	s := &ServiceQuotasServices{
		Table: ui.NewTable([]string{
			"NAME",
			"CODE",
		}, 1, 0),
		sqClient: sqClient,
		app:      app,
	}
	return s
}

func (s ServiceQuotasServices) GetService() string {
	return "Service Quotas"
}

func (s ServiceQuotasServices) GetLabels() []string {
	return []string{}
}

func (s ServiceQuotasServices) viewQuotasHandler() {
	serviceName, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	serviceCode, err := s.GetColSelection("CODE")
	if err != nil {
		return
	}
	quotasView := NewServiceQuotasQuotas(s.sqClient, serviceName, serviceCode, s.app)
	s.app.AddAndSwitch(quotasView)
}

func (s ServiceQuotasServices) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
			Description: "View Quotas",
			Action:      s.viewQuotasHandler,
		},
	}
}

func (s ServiceQuotasServices) Render() {
	pg := sq.NewListServicesPaginator(
		s.sqClient,
		&sq.ListServicesInput{},
	)
	var services []sqTypes.ServiceInfo
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		services = append(services, out.Services...)
	}

	var data [][]string
	for _, v := range services {
		var name, code string
		if v.ServiceName != nil {
			name = *v.ServiceName
		}
		if v.ServiceCode != nil {
			code = *v.ServiceCode
		}
		data = append(data, []string{
			name,
			code,
		})
	}
	s.SetData(data)
}
