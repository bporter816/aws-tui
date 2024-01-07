package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
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

func (c CFFunctions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFFunctions) Render() {
	model, err := c.repo.ListFunctions()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, comment, status, stage, created, modified string
		if v.Name != nil {
			name = *v.Name
		}
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
			name,
			comment,
			status,
			stage,
			created,
			modified,
		})
	}
	c.SetData(data)
}
