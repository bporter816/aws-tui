package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewCFDistributionInvalidations(repo *repo.CloudFront, distributionId string, app *Application) *TableView[model.CloudFrontInvalidation] {
	return NewTableView(
		"CloudFront",
		[]string{distributionId, "Invalidations"},
		[]string{"ID", "STATUS", "CREATED"},
		func(ctx context.Context) ([]model.CloudFrontInvalidation, error) {
			return repo.ListInvalidations(ctx, distributionId)
		},
		func(v model.CloudFrontInvalidation) []string {
			created := "-"
			if v.CreateTime != nil {
				created = v.CreateTime.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.Id, "-"),
				utils.DerefString(v.Status, "-"),
				created,
			}
		},
		func(v *TableView[model.CloudFrontInvalidation]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Paths",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewCFDistributionInvalidationPaths(repo, distributionId, id, app))
					},
				},
			}
		},
	)
}
