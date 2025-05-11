package internal

import (
	"fmt"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheGroups struct {
	*ui.Table
	view.ElastiCache
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

func (e ElastiCacheGroups) GetLabels() []string {
	return []string{"Groups"}
}

func (e ElastiCacheGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].ARN; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var status string
		if v.Status != nil {
			status = utils.TitleCase(*v.Status)
		}
		data = append(data, []string{
			utils.DerefString(v.UserGroupId, ""),
			status,
			fmt.Sprintf("%v", len(v.UserIds)),
			fmt.Sprintf("%v", len(v.ReplicationGroups)),
		})
	}
	e.SetData(data)
}
