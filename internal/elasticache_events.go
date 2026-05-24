package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewElastiCacheEvents(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheEvent] {
	return NewTableView(
		"ElastiCache",
		[]string{"Events"},
		[]string{"DATE", "SOURCE", "TYPE", "MESSAGE"},
		repo.ListEvents,
		func(v model.ElastiCacheEvent) []string {
			var date string
			if v.Date != nil {
				date = v.Date.Format(utils.DefaultTimeFormat)
			}
			return []string{
				date,
				utils.DerefString(v.SourceIdentifier, ""),
				string(v.SourceType),
				utils.DerefString(v.Message, ""),
			}
		},
		nil,
	)
}

func init() {
	register("ElastiCache.Events", func(r Repos, app *Application) Component {
		return NewElastiCacheEvents(r.ElastiCache, app)
	})
}
