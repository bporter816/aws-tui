package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"time"
)

type ElasticacheEvents struct {
	*ui.Table
	ecClient *ec.Client
	app      *Application
}

func NewElasticacheEvents(ecClient *ec.Client, app *Application) *ElasticacheEvents {
	e := &ElasticacheEvents{
		Table: ui.NewTable([]string{
			"DATE",
			"SOURCE",
			"TYPE",
			"MESSAGE",
		}, 1, 0),
		ecClient: ecClient,
		app:      app,
	}
	return e
}

func (e ElasticacheEvents) GetService() string {
	return "Elasticache"
}

func (e ElasticacheEvents) GetLabels() []string {
	return []string{"Events"}
}

func (e ElasticacheEvents) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheEvents) Render() {
	oneWeekAgo := time.Now().AddDate(0, 0, -13) // TODO get this closer to the max 14 days
	pg := ec.NewDescribeEventsPaginator(
		e.ecClient,
		&ec.DescribeEventsInput{
			StartTime: aws.Time(oneWeekAgo),
		},
	)
	var events []ecTypes.Event
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		events = append(events, out.Events...)
	}

	var data [][]string
	for _, v := range events {
		var date, sourceId, sourceType, message string
		if v.Date != nil {
			date = v.Date.Format(utils.DefaultTimeFormat)
		}
		if v.SourceIdentifier != nil {
			sourceId = *v.SourceIdentifier
		}
		sourceType = string(v.SourceType)
		if v.Message != nil {
			message = *v.Message
		}
		data = append(data, []string{
			date,
			sourceId,
			sourceType,
			message,
		})
	}
	e.SetData(data)
}
