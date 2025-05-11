package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type MSKClusters struct {
	*ui.Table
	view.MSK
	repo  *repo.MSK
	app   *Application
	model []model.MSKCluster
}

func NewMSKClusters(repo *repo.MSK, app *Application) *MSKClusters {
	m := &MSKClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"STATE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return m
}

func (m MSKClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (m MSKClusters) tagsHandler() {
	row, err := m.GetRowSelection()
	if err != nil {
		return
	}
	if arn := m.model[row-1].ClusterArn; arn != nil {
		tagsView := NewTags(m.repo, m.GetService(), *arn, m.app)
		m.app.AddAndSwitch(tagsView)
	}
}

func (m MSKClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      m.tagsHandler,
		},
	}
}

func (m *MSKClusters) Render() {
	model, err := m.repo.ListClusters()
	if err != nil {
		panic(err)
	}
	m.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.ClusterName, ""),
			utils.AutoCase(string(v.ClusterType)),
			utils.AutoCase(string(v.State)),
			// TODO add provisioned and serverless details
		})
	}
	m.SetData(data)
}
