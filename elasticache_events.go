package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type ElasticacheEvents struct {
	*ui.Table
	repo *repo.Elasticache
	app  *Application
}

func NewElasticacheEvents(repo *repo.Elasticache, app *Application) *ElasticacheEvents {
	e := &ElasticacheEvents{
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
	model, err := e.repo.ListEvents()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
