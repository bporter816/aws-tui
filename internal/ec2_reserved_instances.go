package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type EC2ReservedInstances struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2ReservedInstances(repo *repo.EC2, app *Application) *EC2ReservedInstances {
	e := &EC2ReservedInstances{
		Table: ui.NewTable([]string{
			"ID",
			"INSTANCE TYPE",
			"SCOPE",
			"STATE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2ReservedInstances) GetLabels() []string {
	return []string{"Reserved Instances"}
}

func (e EC2ReservedInstances) tagsHandler() {
	id, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), id, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2ReservedInstances) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2ReservedInstances) Render() {
	model, err := e.repo.ListReservedInstances(nil)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.ReservedInstancesId, ""),
			string(v.InstanceType),
			string(v.Scope),
			utils.AutoCase(string(v.State)),
		})
	}
	e.SetData(data)
}
