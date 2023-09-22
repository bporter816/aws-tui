package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type ElasticacheGroups struct {
	*ui.Table
	repo *repo.Elasticache
	app  *Application
}

func NewElasticacheGroups(repo *repo.Elasticache, app *Application) *ElasticacheGroups {
	e := &ElasticacheGroups{
		Table: ui.NewTable([]string{
			"ID",
			"STATUS",
			"USERS",
			"CLUSTERS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElasticacheGroups) GetService() string {
	return "Elasticache"
}

func (e ElasticacheGroups) GetLabels() []string {
	return []string{"Groups"}
}

func (e ElasticacheGroups) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheGroups) Render() {
	model, err := e.repo.ListGroups()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, status, users, clusters string
		if v.UserGroupId != nil {
			id = *v.UserGroupId
		}
		if v.Status != nil {
			status = *v.Status
		}
		users = fmt.Sprintf("%v", len(v.UserIds))
		clusters = fmt.Sprintf("%v", len(v.ReplicationGroups))
		data = append(data, []string{
			id,
			status,
			users,
			clusters,
		})
	}
	e.SetData(data)
}
