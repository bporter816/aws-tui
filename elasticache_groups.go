package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type ElasticacheGroups struct {
	*ui.Table
	repo  *repo.Elasticache
	app   *Application
	model []model.ElasticacheGroup
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

func (e ElasticacheGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	if e.model[row-1].ARN == nil {
		return
	}
	tagsView := NewElasticacheTags(e.repo, *e.model[row-1].ARN, name, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ElasticacheGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElasticacheGroups) Render() {
	model, err := e.repo.ListGroups()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var id, status, users, clusters string
		if v.UserGroupId != nil {
			id = *v.UserGroupId
		}
		if v.Status != nil {
			status = utils.TitleCase(*v.Status)
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
