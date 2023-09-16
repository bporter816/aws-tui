package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type SQSQueues struct {
	*ui.Table
	repo *repo.SQS
	app  *Application
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

func (s SQSQueues) GetService() string {
	return "SQS"
}

func (s SQSQueues) GetLabels() []string {
	return []string{"Queues"}
}

func (s SQSQueues) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SQSQueues) Render() {
	model, err := s.repo.ListQueues()
	if err != nil {
		panic(err)
	}
	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Name,
			utils.BoolToString(v.IsFifo, "FIFO", "Standard"),
		})
	}
	s.SetData(data)
}
