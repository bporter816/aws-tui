package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EKSClusters struct {
	*ui.Table
	view.EKS
	repo  *repo.EKS
	app   *Application
	model []model.EKSCluster
}

func NewEKSClusters(repo *repo.EKS, app *Application) *EKSClusters {
	e := &EKSClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"K8S VERSION",
			"CREATED",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EKSClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (e EKSClusters) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].Arn; arn != nil {
		tagsView := NewEKSTags(e.repo, *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e EKSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *EKSClusters) Render() {
	model, err := e.repo.ListClusters()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var name, status, k8sVersion, created string
		if v.Name != nil {
			name = *v.Name
		}
		status = utils.TitleCase(string(v.Status))
		if v.Version != nil {
			k8sVersion = *v.Version
		}
		if v.CreatedAt != nil {
			created = v.CreatedAt.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			name,
			status,
			k8sVersion,
			created,
		})
	}
	e.SetData(data)
}
