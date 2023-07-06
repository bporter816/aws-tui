package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
)

type ElasticacheReservedCacheNodes struct {
	*Table
	ecClient *ec.Client
	app      *Application
}

func NewElasticacheReservedCacheNodes(ecClient *ec.Client, app *Application) *ElasticacheReservedCacheNodes {
	e := &ElasticacheReservedCacheNodes{
		Table: NewTable([]string{
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

func (e ElasticacheReservedCacheNodes) GetName() string {
	return "Elasticache | Reserved Nodes"
}

func (e ElasticacheReservedCacheNodes) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheReservedCacheNodes) Render() {
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
	for _, v := range reservations {
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