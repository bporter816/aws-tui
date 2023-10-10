package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ElastiCacheReservedCacheNodes struct {
	*ui.Table
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

func (e ElastiCacheReservedCacheNodes) GetService() string {
	return "ElastiCache"
}

func (e ElastiCacheReservedCacheNodes) GetLabels() []string {
	return []string{"Reserved Nodes"}
}

func (e ElastiCacheReservedCacheNodes) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	if e.model[row-1].ReservationARN == nil {
		return
	}
	tagsView := NewElastiCacheTags(e.repo, *e.model[row-1].ReservationARN, name, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ElastiCacheReservedCacheNodes) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
			*v.ReservedCacheNodeId,
			*v.OfferingType,
			utils.TitleCase(*v.ProductDescription),
			*v.CacheNodeType,
			strconv.Itoa(int(v.CacheNodeCount)),
			utils.TitleCase(*v.State),
		})
	}
	e.SetData(data)
}
