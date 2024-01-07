package main

import (
	"fmt"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EC2Subnets struct {
	*ui.Table
	view.EC2
	repo      *repo.EC2
	subnetIds []string
	label     string
	app       *Application
}

func NewEC2Subnets(repo *repo.EC2, subnetIds []string, label string, app *Application) *EC2Subnets {
	e := &EC2Subnets{
		Table: ui.NewTable([]string{
			"NAME",
			"SUBNET ID",
			"STATE",
			"AVAILABILITY ZONE",
			"IPV4 CIDR",
			"VPC ID",
		}, 1, 0),
		repo:      repo,
		subnetIds: subnetIds,
		label:     label,
		app:       app,
	}
	return e
}

func (e EC2Subnets) GetLabels() []string {
	if len(e.subnetIds) > 0 {
		return []string{e.label, "Subnets"}
	} else {
		return []string{"Subnets"}
	}
}

func (e EC2Subnets) tagsHandler() {
	subnetId, err := e.GetColSelection("SUBNET ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SubnetTags(e.repo, subnetId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2Subnets) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2Subnets) Render() {
	model, err := e.repo.ListSubnets(e.subnetIds)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, subnetId, state, availabilityZone, ipv4Cidr, vpcId string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		} else {
			name = "-"
		}
		if v.SubnetId != nil {
			subnetId = *v.SubnetId
		}
		state = utils.TitleCase(string(v.State))
		if v.AvailabilityZone != nil {
			availabilityZone = *v.AvailabilityZone
			if v.AvailabilityZoneId != nil {
				availabilityZone += fmt.Sprintf(" (%v)", *v.AvailabilityZoneId)
			}
		}
		if v.CidrBlock != nil {
			ipv4Cidr = *v.CidrBlock
		}
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		data = append(data, []string{
			name,
			subnetId,
			state,
			availabilityZone,
			ipv4Cidr,
			vpcId,
		})
	}
	e.SetData(data)
}
