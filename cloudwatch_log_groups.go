package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"strconv"
)

type CloudwatchLogGroups struct {
	*ui.Table
	repo *repo.Cloudwatch
	app  *Application
}

func NewCloudwatchLogGroups(repo *repo.Cloudwatch, app *Application) *CloudwatchLogGroups {
	c := &CloudwatchLogGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"RETENTION",
			"DATA PROTECTION",
			"METRIC FILTERS",
			"STORED DATA",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return c
}

func (c CloudwatchLogGroups) GetService() string {
	return "Cloudwatch"
}

func (c CloudwatchLogGroups) GetLabels() []string {
	return []string{"Log Groups"}
}

func (c CloudwatchLogGroups) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CloudwatchLogGroups) Render() {
	model, err := c.repo.ListLogGroups()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, retention, dataProtection, metricFilters, storedData string
		if v.LogGroupName != nil {
			name = *v.LogGroupName
		}
		if v.RetentionInDays != nil {
			// TODO have nice output for months, years
			retention = strconv.Itoa(int(*v.RetentionInDays))
		} else {
			retention = "Never expire"
		}
		dataProtection = utils.TitleCase(string(v.DataProtectionStatus))
		if len(dataProtection) == 0 {
			dataProtection = "-"
		}
		if v.MetricFilterCount != nil {
			metricFilters = strconv.Itoa(int(*v.MetricFilterCount))
		}
		if v.StoredBytes != nil {
			// TODO print kb, mb, etc
			storedData = strconv.FormatInt(*v.StoredBytes, 10) + " B"
		}
		data = append(data, []string{
			name,
			retention,
			dataProtection,
			metricFilters,
			storedData,
		})
	}
	c.SetData(data)
}
