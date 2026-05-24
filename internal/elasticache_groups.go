package internal

import (
	"fmt"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheGroups(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheGroup] {
	return NewTableView(
		"ElastiCache",
		[]string{"Groups"},
		[]string{"ID", "STATUS", "USERS", "CLUSTERS"},
		repo.ListGroups,
		func(v model.ElastiCacheGroup) []string {
			var status string
			if v.Status != nil {
				status = utils.TitleCase(*v.Status)
			}
			return []string{
				utils.DerefString(v.UserGroupId, ""),
				status,
				fmt.Sprintf("%v", len(v.UserIds)),
				fmt.Sprintf("%v", len(v.ReplicationGroups)),
			}
		},
		func(v *TableView[model.ElastiCacheGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ARN == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ARN, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ElastiCache.Groups", func(r Repos, app *Application) Component {
		return NewElastiCacheGroups(r.ElastiCache, app)
	})
}
