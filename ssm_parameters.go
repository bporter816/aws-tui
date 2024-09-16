package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type SSMParameters struct {
	*ui.Table
	view.SSM
	repo  *repo.SSM
	app   *Application
	model []model.SSMParameter
}

func NewSSMParameters(repo *repo.SSM, app *Application) *SSMParameters {
	s := &SSMParameters{
		Table: ui.NewTable([]string{
			"NAME",
			"TIER",
			"TYPE",
			"DATA TYPE",
			"VERSIONS",
			"POLICIES",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return s
}

func (s SSMParameters) GetLabels() []string {
	return []string{"Parameters"}
}

func (s SSMParameters) tagsHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	if arn := s.model[row-1].Name; arn != nil {
		tagsView := NewTags(s.repo, s.GetService(), "parameter:"+*arn, s.app)
		s.app.AddAndSwitch(tagsView)
	}
}

func (s SSMParameters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      s.tagsHandler,
		},
	}
}

func (s *SSMParameters) Render() {
	model, err := s.repo.ListParameters()
	if err != nil {
		panic(err)
	}
	s.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			string(v.Tier),
			string(v.Type),
			utils.DerefString(v.DataType, ""),
			strconv.FormatInt(v.Version, 10),
			strconv.Itoa(len(v.Policies)),
		})
	}
	s.SetData(data)
}
