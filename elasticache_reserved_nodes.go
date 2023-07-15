package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
)

type ElasticacheReservedCacheNodes struct {
	*ui.Table
	ecClient *ec.Client
	app      *Application
	arns     []string
}

func NewElasticacheReservedCacheNodes(ecClient *ec.Client, app *Application) *ElasticacheReservedCacheNodes {
	e := &ElasticacheReservedCacheNodes{
		Table: ui.NewTable([]string{
			"ID",
			"OFFERING TYPE",
			"ENGINE",
			"NODE TYPE",
			"NODES",
			"STATUS",
		}, 1, 0),
		ecClient: ecClient,
		app:      app,
	}
	return e
}

func (e ElasticacheReservedCacheNodes) GetService() string {
	return "Elasticache"
}

func (e ElasticacheReservedCacheNodes) GetLabels() []string {
	return []string{"Reserved Nodes"}
}

func (e ElasticacheReservedCacheNodes) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewElasticacheTags(e.ecClient, ElasticacheResourceTypeReservedNode, e.arns[row-1], name, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ElasticacheReservedCacheNodes) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElasticacheReservedCacheNodes) Render() {
	pg := ec.NewDescribeReservedCacheNodesPaginator(
		e.ecClient,
		&ec.DescribeReservedCacheNodesInput{},
	)
	var reservations []ecTypes.ReservedCacheNode
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		reservations = append(reservations, out.ReservedCacheNodes...)
	}

	caser := cases.Title(language.English)
	var data [][]string
	e.arns = make([]string, len(reservations))
	for i, v := range reservations {
		e.arns[i] = *v.ReservationARN
		data = append(data, []string{
			*v.ReservedCacheNodeId,
			*v.OfferingType,
			caser.String(*v.ProductDescription),
			*v.CacheNodeType,
			strconv.Itoa(int(v.CacheNodeCount)),
			caser.String(*v.State),
		})
	}
	e.SetData(data)
}
