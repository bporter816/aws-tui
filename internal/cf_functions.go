package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewCFFunctions(repo *repo.CloudFront, app *Application) *TableView[model.CloudFrontFunction] {
	return NewTableView(
		"CloudFront",
		[]string{"Functions"},
		[]string{"NAME", "COMMENT", "STATUS", "STAGE", "CREATED", "MODIFIED"},
		repo.ListFunctions,
		func(v model.CloudFrontFunction) []string {
			var comment, status, stage, created, modified string
			if v.FunctionConfig != nil && v.FunctionConfig.Comment != nil {
				comment = *v.FunctionConfig.Comment
			}
			if v.Status != nil {
				status = utils.TitleCase(*v.Status)
			}
			if v.FunctionMetadata != nil {
				stage = utils.TitleCase(string(v.FunctionMetadata.Stage))
				if v.FunctionMetadata.CreatedTime != nil {
					created = v.FunctionMetadata.CreatedTime.Format(utils.DefaultTimeFormat)
				}
				if v.FunctionMetadata.LastModifiedTime != nil {
					modified = v.FunctionMetadata.LastModifiedTime.Format(utils.DefaultTimeFormat)
				}
			}
			return []string{
				utils.DerefString(v.Name, ""),
				comment,
				status,
				stage,
				created,
				modified,
			}
		},
		func(v *TableView[model.CloudFrontFunction]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
					Description: "View Code",
					Action: func() {
						name, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						stage, err := v.GetColSelection("STAGE")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewCFFunctionCode(repo, name, stage, app))
					},
				},
			}
		},
	)
}

func init() {
	register("CloudFront.Functions", func(r Repos, app *Application) Component {
		return NewCFFunctions(r.CloudFront, app)
	})
}
