package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type SQSQueues struct {
	*ui.Table
	view.SQS
	repo  *repo.SQS
	app   *Application
	model []model.SQSQueue
}

func NewSQSQueues(repo *repo.SQS, app *Application) *SQSQueues {
	s := &SQSQueues{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return s
}

func (s SQSQueues) GetLabels() []string {
	return []string{"Queues"}
}

func (s SQSQueues) accessPolicyHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	accessPolicyView := NewSQSAccessPolicy(s.repo, s.model[row-1].QueueUrl, s.app)
	s.app.AddAndSwitch(accessPolicyView)
}

func (s SQSQueues) tagsHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	tagsView := NewTags(s.repo, s.GetService(), s.model[row-1].QueueUrl, s.app)
	s.app.AddAndSwitch(tagsView)
}

func (s SQSQueues) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Access Policy",
			Action:      s.accessPolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      s.tagsHandler,
		},
	}
}

func (s *SQSQueues) Render() {
	model, err := s.repo.ListQueues()
	if err != nil {
		panic(err)
	}
	s.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Name,
			utils.BoolToString(v.IsFifo, "FIFO", "Standard"),
		})
	}
	s.SetData(data)
}
