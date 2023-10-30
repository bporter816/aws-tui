package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheGroups struct {
	*ui.Table
	repo  *repo.ElastiCache
	app   *Application
	model []model.ElastiCacheGroup
}

func NewElastiCacheGroups(repo *repo.ElastiCache, app *Application) *ElastiCacheGroups {
	e := &ElastiCacheGroups{
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

func (e ElastiCacheGroups) GetService() string {
	return "ElastiCache"
}

func (e ElastiCacheGroups) GetLabels() []string {
	return []string{"Groups"}
}

func (e ElastiCacheGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	if arn := e.model[row-1].ARN; arn != nil {
		tagsView := NewElastiCacheTags(e.repo, *arn, name, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElastiCacheGroups) Render() {
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
