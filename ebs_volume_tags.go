package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EBSVolumeTags struct {
	*ui.Table
	view.EBS
	volumeId string
	repo     *repo.EC2
	app      *Application
}

func NewEBSVolumeTags(volumeId string, repo *repo.EC2, app *Application) *EBSVolumeTags {
	e := &EBSVolumeTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		volumeId: volumeId,
		repo:     repo,
		app:      app,
	}
	return e
}

func (e EBSVolumeTags) GetLabels() []string {
	return []string{e.volumeId, "Tags"}
}

func (e EBSVolumeTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EBSVolumeTags) Render() {
	model, err := e.repo.ListTags(e.volumeId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	e.SetData(data)
}
