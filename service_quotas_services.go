package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type ServiceQuotasServices struct {
	*ui.Table
	repo *repo.ServiceQuotas
	app  *Application
}

func NewServiceQuotasServices(repo *repo.ServiceQuotas, app *Application) *ServiceQuotasServices {
	s := &ServiceQuotasServices{
		Table: ui.NewTable([]string{
			"NAME",
			"CODE",
		}, 1, 0),
		repo: repo,
		app:  app,
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
	quotasView := NewServiceQuotasQuotas(s.repo, serviceName, serviceCode, s.app)
	s.app.AddAndSwitch(quotasView)
}

func (s ServiceQuotasServices) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
			Description: "View Quotas",
			Action:      s.viewQuotasHandler,
		},
	}
}

func (s ServiceQuotasServices) Render() {
	model, err := s.repo.ListServices()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
