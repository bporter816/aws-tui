package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ServiceQuotasServices struct {
	*ui.Table
	view.ServiceQuotas
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
		data = append(data, []string{
			utils.DerefString(v.ServiceName, ""),
			utils.DerefString(v.ServiceCode, ""),
		})
	}
	s.SetData(data)
}
