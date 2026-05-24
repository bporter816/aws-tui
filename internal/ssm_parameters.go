package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewSSMParameters(repo *repo.SSM, app *Application) *TableView[model.SSMParameter] {
	return NewTableView(
		"Systems Manager",
		[]string{"Parameters"},
		[]string{"NAME", "TIER", "TYPE", "DATA TYPE", "VERSIONS", "POLICIES"},
		repo.ListParameters,
		func(v model.SSMParameter) []string {
			return []string{
				utils.DerefString(v.Name, ""),
				string(v.Tier),
				string(v.Type),
				utils.DerefString(v.DataType, ""),
				strconv.FormatInt(v.Version, 10),
				strconv.Itoa(len(v.Policies)),
			}
		},
		func(v *TableView[model.SSMParameter]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.Name == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), "parameter:"+*item.Name, app))
					},
				},
			}
		},
	)
}
