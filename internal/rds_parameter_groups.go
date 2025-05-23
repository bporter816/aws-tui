package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type RDSParameterGroups struct {
	*ui.Table
	view.RDS
	repo  *repo.RDS
	app   *Application
	model []model.ModelWithArn
}

func NewRDSParameterGroups(repo *repo.RDS, app *Application) *RDSParameterGroups {
	r := &RDSParameterGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"FAMILY",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return r
}

func (r RDSParameterGroups) GetLabels() []string {
	return []string{"Parameter Groups"}
}

func (r RDSParameterGroups) parametersHandler() {
	name, err := r.GetColSelection("NAME")
	if err != nil {
		return
	}
	groupType, err := r.GetColSelection("TYPE")
	if err != nil {
		return
	}
	var parameterGroupType model.RDSParameterGroupType
	if groupType == "Cluster" {
		parameterGroupType = model.RDSParameterGroupTypeCluster
	} else {
		parameterGroupType = model.RDSParameterGroupTypeInstance
	}
	parametersView := NewRDSParameters(r.repo, r.app, name, parameterGroupType)
	r.app.AddAndSwitch(parametersView)
}

func (r RDSParameterGroups) tagsHandler() {
	row, err := r.GetRowSelection()
	if err != nil {
		return
	}
	tagsView := NewTags(r.repo, r.GetService(), r.model[row-1].Arn(), r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r RDSParameterGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Parameters",
			Action:      r.parametersHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r *RDSParameterGroups) Render() {
	clusterParameterGroups, err := r.repo.ListClusterParameterGroups()
	if err != nil {
		panic(err)
	}

	instanceParameterGroups, err := r.repo.ListInstanceParameterGroups()
	if err != nil {
		panic(err)
	}
	var objects []model.ModelWithArn

	var data [][]string
	for _, v := range clusterParameterGroups {
		objects = append(objects, v)
		data = append(data, []string{
			utils.DerefString(v.DBClusterParameterGroupName, ""),
			"Cluster",
			utils.DerefString(v.DBParameterGroupFamily, ""),
			utils.DerefString(v.Description, ""),
		})
	}
	for _, v := range instanceParameterGroups {
		objects = append(objects, v)
		data = append(data, []string{
			utils.DerefString(v.DBParameterGroupName, ""),
			"Instance",
			utils.DerefString(v.DBParameterGroupFamily, ""),
			utils.DerefString(v.Description, ""),
		})
	}
	r.model = objects
	r.SetData(data)
}
