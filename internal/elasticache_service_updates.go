package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheServiceUpdates(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheServiceUpdate] {
	return NewTableView(
		"ElastiCache",
		[]string{"Service Updates"},
		[]string{"NAME", "ENGINE VERSION", "TYPE", "SEVERITY", "STATUS", "RELEASED", "APPLY BY", "AUTO APPLY"},
		repo.ListServiceUpdates,
		func(v model.ElastiCacheServiceUpdate) []string {
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
			return []string{
				utils.DerefString(v.ServiceUpdateName, ""),
				utils.DerefString(v.EngineVersion, ""),
				utils.AutoCase(string(v.ServiceUpdateType)),
				utils.AutoCase(string(v.ServiceUpdateSeverity)),
				utils.AutoCase(string(v.ServiceUpdateStatus)),
				released,
				applyBy,
				autoApply,
			}
		},
		func(v *TableView[model.ElastiCacheServiceUpdate]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
					Description: "Status",
					Action: func() {
						serviceUpdateName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewElastiCacheUpdateActions(repo, app, []string{}, []string{}, serviceUpdateName))
					},
				},
			}
		},
	)
}

func init() {
	register("ElastiCache.Service Updates", func(r Repos, app *Application) Component {
		return NewElastiCacheServiceUpdates(r.ElastiCache, app)
	})
}
