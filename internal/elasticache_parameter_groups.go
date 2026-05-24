package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheParameterGroups(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheParameterGroup] {
	return NewTableView(
		"ElastiCache",
		[]string{"Parameter Groups"},
		[]string{"NAME", "FAMILY", "DESCRIPTION", "GLOBAL"},
		repo.ListParameterGroups,
		func(v model.ElastiCacheParameterGroup) []string {
			var isGlobal string
			if v.IsGlobal != nil {
				isGlobal = utils.BoolToString(*v.IsGlobal, "Yes", "No")
			}
			return []string{
				utils.DerefString(v.CacheParameterGroupName, ""),
				utils.DerefString(v.CacheParameterGroupFamily, ""),
				utils.DerefString(v.Description, ""),
				isGlobal,
			}
		},
		func(v *TableView[model.ElastiCacheParameterGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "View Parameters",
					Action: func() {
						name, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewElastiCacheParameters(repo, name, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ElastiCache.Parameter Groups", func(r Repos, app *Application) Component {
		return NewElastiCacheParameterGroups(r.ElastiCache, app)
	})
}
