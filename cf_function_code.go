package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type CFFunctionCode struct {
	*ui.Text
	view.CloudFront
	repo  *repo.CloudFront
	name  string
	stage string
}

func NewCFFunctionCode(repo *repo.CloudFront, name string, stage string, app *Application) *CFFunctionCode {
	c := &CFFunctionCode{
		Text:  ui.NewText(true, "js"), // CloudFront functions can only be JavaScript
		repo:  repo,
		name:  name,
		stage: stage,
	}
	return c
}

func (c CFFunctionCode) GetLabels() []string {
	return []string{c.name + utils.AutoCase(string(c.stage)), "Code"}
}

func (c CFFunctionCode) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFFunctionCode) Render() {
	code, err := c.repo.GetFunctionCode(c.name, c.stage)
	if err != nil {
		panic(err)
	}

	c.SetText(code)
}
