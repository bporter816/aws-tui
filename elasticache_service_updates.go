package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type ElasticacheServiceUpdates struct {
	*ui.Table
	repo *repo.Elasticache
	app  *Application
}

func NewElasticacheServiceUpdates(repo *repo.Elasticache, app *Application) *ElasticacheServiceUpdates {
	e := &ElasticacheServiceUpdates{
		Table: ui.NewTable([]string{
			"NAME",
			"ENGINE VERSION",
			"TYPE",
			"SEVERITY",
			"STATUS",
			"RELEASED",
			"APPLY BY",
			"AUTO APPLY",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElasticacheServiceUpdates) GetService() string {
	return "Elasticache"
}

func (e ElasticacheServiceUpdates) GetLabels() []string {
	return []string{"Service Updates"}
}

func (e ElasticacheServiceUpdates) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheServiceUpdates) Render() {
	model, err := e.repo.ListServiceUpdates()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, engineVersion, serviceUpdateType, severity, status, released, applyBy, autoApply string
		if v.ServiceUpdateName != nil {
			name = *v.ServiceUpdateName
		}
		if v.EngineVersion != nil {
			engineVersion = *v.EngineVersion
		}
		serviceUpdateType = utils.AutoCase(string(v.ServiceUpdateType))
		severity = utils.AutoCase(string(v.ServiceUpdateSeverity))
		status = utils.AutoCase(string(v.ServiceUpdateStatus))
		if v.ServiceUpdateReleaseDate != nil {
			released = v.ServiceUpdateReleaseDate.Format(utils.DefaultTimeFormat)
		}
		if v.ServiceUpdateRecommendedApplyByDate != nil {
			applyBy = v.ServiceUpdateRecommendedApplyByDate.Format(utils.DefaultTimeFormat)
		}
		if v.AutoUpdateAfterRecommendedApplyByDate != nil {
			autoApply = utils.BoolToString(*v.AutoUpdateAfterRecommendedApplyByDate, "Yes", "No")
		}
		data = append(data, []string{
			name,
			engineVersion,
			serviceUpdateType,
			severity,
			status,
			released,
			applyBy,
			autoApply,
		})
	}
	e.SetData(data)
}
