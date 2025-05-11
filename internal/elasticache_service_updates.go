package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheServiceUpdates struct {
	*ui.Table
	view.ElastiCache
	repo *repo.ElastiCache
	app  *Application
}

func NewElastiCacheServiceUpdates(repo *repo.ElastiCache, app *Application) *ElastiCacheServiceUpdates {
	e := &ElastiCacheServiceUpdates{
		Table: ui.NewTable([]string{
			"NAME",
			"ENGINE VERSION",
			"TYPE",
			"SEVERITY",
			"STATUS",
			"RELEASED",
			"APPLY BY",
			"AUTO APPLY",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElastiCacheServiceUpdates) GetLabels() []string {
	return []string{"Service Updates"}
}

func (e ElastiCacheServiceUpdates) statusHandler() {
	serviceUpdateName, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	statusView := NewElastiCacheUpdateActions(e.repo, e.app, []string{}, []string{}, serviceUpdateName)
	e.app.AddAndSwitch(statusView)
}

func (e ElastiCacheServiceUpdates) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			Description: "Status",
			Action:      e.statusHandler,
		},
	}
}

func (e ElastiCacheServiceUpdates) Render() {
	model, err := e.repo.ListServiceUpdates()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var released, applyBy, autoApply string
		if v.ServiceUpdateReleaseDate != nil {
			released = v.ServiceUpdateReleaseDate.Format(utils.DefaultTimeFormat)
		}
		if v.ServiceUpdateRecommendedApplyByDate != nil {
			applyBy = v.ServiceUpdateRecommendedApplyByDate.Format(utils.DefaultTimeFormat)
		}
		if v.AutoUpdateAfterRecommendedApplyByDate != nil {
			autoApply = utils.BoolToString(*v.AutoUpdateAfterRecommendedApplyByDate, "Yes", "No")
		}
		data = append(data, []string{
			utils.DerefString(v.ServiceUpdateName, ""),
			utils.DerefString(v.EngineVersion, ""),
			utils.AutoCase(string(v.ServiceUpdateType)),
			utils.AutoCase(string(v.ServiceUpdateSeverity)),
			utils.AutoCase(string(v.ServiceUpdateStatus)),
			released,
			applyBy,
			autoApply,
		})
	}
	e.SetData(data)
}
