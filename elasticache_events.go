package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"time"
)

type ElasticacheEvents struct {
	*Table
	ecClient *ec.Client
	app      *Application
}

func NewElasticacheEvents(ecClient *ec.Client, app *Application) *ElasticacheEvents {
	e := &ElasticacheEvents{
		Table: NewTable([]string{
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

func (e ElasticacheEvents) GetName() string {
	return "Elasticache | Events"
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
			date = v.Date.Format("2006-01-02 15:04:05")
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
