package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheReservedCacheNodes(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheReservedNode] {
	return NewTableView(
		"ElastiCache",
		[]string{"Reserved Nodes"},
		[]string{"ID", "OFFERING TYPE", "ENGINE", "NODE TYPE", "NODES", "STATUS"},
		repo.ListReservedNodes,
		func(v model.ElastiCacheReservedNode) []string {
			return []string{
				utils.DerefString(v.ReservedCacheNodeId, ""),
				utils.DerefString(v.OfferingType, ""),
				utils.TitleCase(*v.ProductDescription),
				utils.DerefString(v.CacheNodeType, ""),
				strconv.Itoa(int(*v.CacheNodeCount)),
				utils.TitleCase(*v.State),
			}
		},
		func(v *TableView[model.ElastiCacheReservedNode]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ReservationARN == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ReservationARN, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ElastiCache.Reserved Nodes", func(r Repos, app *Application) Component {
		return NewElastiCacheReservedCacheNodes(r.ElastiCache, app)
	})
}
