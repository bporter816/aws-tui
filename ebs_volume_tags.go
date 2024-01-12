package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
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
	model, err := e.repo.ListVolumes(
		[]ec2Types.Filter{
			{
				Name:   aws.String("volume-id"),
				Values: []string{e.volumeId},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	if len(model) != 1 {
		panic("expected exactly one volume spec")
	}

	var data [][]string
	for _, v := range model[0].Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
