package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type ElastiCacheEvents struct {
	*ui.Table
	view.ElastiCache
	repo *repo.ElastiCache
	app  *Application
}

func NewElastiCacheEvents(repo *repo.ElastiCache, app *Application) *ElastiCacheEvents {
	e := &ElastiCacheEvents{
		Table: ui.NewTable([]string{
			"DATE",
			"SOURCE",
			"TYPE",
			"MESSAGE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElastiCacheEvents) GetLabels() []string {
	return []string{"Events"}
}

func (e ElastiCacheEvents) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElastiCacheEvents) Render() {
	model, err := e.repo.ListEvents()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var date string
		if v.Date != nil {
			date = v.Date.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			date,
			utils.DerefString(v.SourceIdentifier, ""),
			string(v.SourceType),
			utils.DerefString(v.Message, ""),
		})
	}
	e.SetData(data)
}
