package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type EKSClusters struct {
	*ui.Table
	repo *repo.EKS
	app  *Application
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

func (e EKSClusters) GetService() string {
	return "EKS"
}

func (e EKSClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (e EKSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EKSClusters) Render() {
	model, err := e.repo.ListClusters()
	if err != nil {
		panic(err)
	}

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
