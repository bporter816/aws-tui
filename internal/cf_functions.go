package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type CFFunctions struct {
	*ui.Table
	view.CloudFront
	repo *repo.CloudFront
	app  *Application
}

func NewCFFunctions(repo *repo.CloudFront, app *Application) *CFFunctions {
	c := &CFFunctions{
		Table: ui.NewTable([]string{
			"NAME",
			"COMMENT",
			"STATUS",
			"STAGE",
			"CREATED",
			"MODIFIED",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return c
}

func (c CFFunctions) GetLabels() []string {
	return []string{"Functions"}
}

func (c CFFunctions) codeHandler() {
	name, err := c.GetColSelection("NAME")
	if err != nil {
		return
	}
	stage, err := c.GetColSelection("STAGE")
	if err != nil {
		return
	}
	codeView := NewCFFunctionCode(c.repo, name, stage, c.app)
	c.app.AddAndSwitch(codeView)
}

func (c CFFunctions) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "View Code",
			Action:      c.codeHandler,
		},
	}
}

func (c CFFunctions) Render() {
	model, err := c.repo.ListFunctions()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			comment,
			status,
			stage,
			created,
			modified,
		})
	}
	c.SetData(data)
}
