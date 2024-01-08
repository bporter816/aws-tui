package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EC2ReservedInstanceTags struct {
	*ui.Table
	view.EC2
	reservedInstanceId string
	repo               *repo.EC2
	app                *Application
}

func NewEC2ReservedInstanceTags(reservedInstanceId string, repo *repo.EC2, app *Application) *EC2ReservedInstanceTags {
	e := &EC2ReservedInstanceTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		reservedInstanceId: reservedInstanceId,
		repo:               repo,
		app:                app,
	}
	return e
}

func (e EC2ReservedInstanceTags) GetLabels() []string {
	return []string{e.reservedInstanceId, "Tags"}
}

func (e EC2ReservedInstanceTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2ReservedInstanceTags) Render() {
	model, err := e.repo.ListReservedInstances(
		[]ec2Types.Filter{
			{
				Name:   aws.String("reserved-instances-id"),
				Values: []string{e.reservedInstanceId},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	if len(model) != 1 {
		panic("expected exactly one reserved instances spec")
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
