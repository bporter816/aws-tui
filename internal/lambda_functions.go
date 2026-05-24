package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type lambdaRepo interface {
	ListFunctions() ([]model.LambdaFunction, error)
	ListTags(string) (model.Tags, error)
}

func NewLambdaFunctions(repo lambdaRepo, app *Application) *TableView[model.LambdaFunction] {
	return NewTableView(
		"Lambda",
		[]string{"Functions"},
		[]string{"NAME", "STATUS", "RUNTIME", "CODE SIZE", "MEMORY", "TIMEOUT", "LAYERS", "DESCRIPTION"},
		repo.ListFunctions,
		func(f model.LambdaFunction) []string {
			var memory, timeout string
			if f.MemorySize != nil {
				memory = utils.FormatSize(int64(*f.MemorySize)<<20, 1)
			}
			if f.Timeout != nil {
				timeout = strconv.FormatInt(int64(*f.Timeout), 10) + " sec"
			}
			return []string{
				utils.DerefString(f.FunctionName, ""),
				utils.TitleCase(string(f.State)),
				string(f.Runtime),
				utils.FormatSize(f.CodeSize, 1),
				memory,
				timeout,
				strconv.Itoa(len(f.Layers)),
				utils.DerefString(f.Description, ""),
			}
		},
		func(v *TableView[model.LambdaFunction]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						f, err := v.Selected()
						if err != nil || f.FunctionArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *f.FunctionArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("Lambda.Functions", func(r Repos, app *Application) Component {
		return NewLambdaFunctions(r.Lambda, app)
	})
}
