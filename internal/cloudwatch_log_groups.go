package internal

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewCloudWatchLogGroups(repo *repo.CloudWatch, app *Application) *TableView[model.CloudWatchLogGroup] {
	return NewTableView(
		"CloudWatch",
		[]string{"Log Groups"},
		[]string{"NAME", "RETENTION", "DATA PROTECTION", "METRIC FILTERS", "STORED DATA"},
		repo.ListLogGroups,
		func(v model.CloudWatchLogGroup) []string {
			var retention, metricFilters, storedData string
			if v.RetentionInDays != nil {
				// TODO have nice output for months, years
				retention = strconv.Itoa(int(*v.RetentionInDays))
			} else {
				retention = "Never expire"
			}
			dataProtection := utils.TitleCase(string(v.DataProtectionStatus))
			if len(dataProtection) == 0 {
				dataProtection = "-"
			}
			if v.MetricFilterCount != nil {
				metricFilters = strconv.Itoa(int(*v.MetricFilterCount))
			}
			if v.StoredBytes != nil {
				storedData = utils.FormatSize(*v.StoredBytes, 1)
			}
			return []string{
				utils.DerefString(v.LogGroupName, ""),
				retention,
				dataProtection,
				metricFilters,
				storedData,
			}
		},
		func(v *TableView[model.CloudWatchLogGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.Arn == nil {
							return
						}
						arnStr := strings.TrimSuffix(*item.Arn, ":*")
						app.AddAndSwitch(NewTags(repo, v.GetService(), arnStr, app))
					},
				},
			}
		},
	)
}

func init() {
	register("CloudWatch.Log Groups", func(r Repos, app *Application) Component {
		return NewCloudWatchLogGroups(r.CloudWatch, app)
	})
}
