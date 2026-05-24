package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewServiceQuotasServices(repo *repo.ServiceQuotas, app *Application) *TableView[model.ServiceQuotasService] {
	return NewTableView(
		"Service Quotas",
		[]string{},
		[]string{"NAME", "CODE"},
		repo.ListServices,
		func(v model.ServiceQuotasService) []string {
			return []string{
				utils.DerefString(v.ServiceName, ""),
				utils.DerefString(v.ServiceCode, ""),
			}
		},
		func(v *TableView[model.ServiceQuotasService]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
					Description: "View Quotas",
					Action: func() {
						serviceName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						serviceCode, err := v.GetColSelection("CODE")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewServiceQuotasQuotas(repo, serviceName, serviceCode, app))
					},
				},
			}
		},
	)
}
