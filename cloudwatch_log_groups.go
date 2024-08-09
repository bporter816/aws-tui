package main

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type CloudWatchLogGroups struct {
	*ui.Table
	view.CloudWatch
	repo  *repo.CloudWatch
	app   *Application
	model []model.CloudWatchLogGroup
}

func NewCloudWatchLogGroups(repo *repo.CloudWatch, app *Application) *CloudWatchLogGroups {
	c := &CloudWatchLogGroups{
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

func (c CloudWatchLogGroups) GetLabels() []string {
	return []string{"Log Groups"}
}

func (c CloudWatchLogGroups) tagsHandler() {
	row, err := c.GetRowSelection()
	if err != nil {
		return
	}
	if c.model[row-1].Arn == nil {
		return
	}
	arn := *c.model[row-1].Arn
	if strings.HasSuffix(arn, ":*") {
		arn = arn[:len(arn)-2]
	}
	tagsView := NewTags(c.repo, c.GetService(), arn, c.app)
	c.app.AddAndSwitch(tagsView)
}

func (c CloudWatchLogGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      c.tagsHandler,
		},
	}
}

func (c *CloudWatchLogGroups) Render() {
	model, err := c.repo.ListLogGroups()
	if err != nil {
		panic(err)
	}
	c.model = model

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
			storedData = utils.FormatSize(*v.StoredBytes, 1)
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
