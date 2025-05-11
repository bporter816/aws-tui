package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheReservedCacheNodes struct {
	*ui.Table
	view.ElastiCache
	repo  *repo.ElastiCache
	app   *Application
	model []model.ElastiCacheReservedNode
}

func NewElastiCacheReservedCacheNodes(repo *repo.ElastiCache, app *Application) *ElastiCacheReservedCacheNodes {
	e := &ElastiCacheReservedCacheNodes{
		Table: ui.NewTable([]string{
			"ID",
			"OFFERING TYPE",
			"ENGINE",
			"NODE TYPE",
			"NODES",
			"STATUS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElastiCacheReservedCacheNodes) GetLabels() []string {
	return []string{"Reserved Nodes"}
}

func (e ElastiCacheReservedCacheNodes) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].ReservationARN; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheReservedCacheNodes) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElastiCacheReservedCacheNodes) Render() {
	model, err := e.repo.ListReservedNodes()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.ReservedCacheNodeId, ""),
			utils.DerefString(v.OfferingType, ""),
			utils.TitleCase(*v.ProductDescription),
			utils.DerefString(v.CacheNodeType, ""),
			strconv.Itoa(int(*v.CacheNodeCount)),
			utils.TitleCase(*v.State),
		})
	}
	e.SetData(data)
}
