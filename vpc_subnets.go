package main

import (
	"fmt"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type VPCSubnets struct {
	*ui.Table
	view.VPC
	repo      *repo.EC2
	subnetIds []string
	label     string
	app       *Application
}

func NewVPCSubnets(repo *repo.EC2, subnetIds []string, label string, app *Application) *VPCSubnets {
	e := &VPCSubnets{
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

func (e VPCSubnets) GetLabels() []string {
	if len(e.subnetIds) > 0 {
		return []string{e.label, "Subnets"}
	} else {
		return []string{"Subnets"}
	}
}

func (e VPCSubnets) tagsHandler() {
	subnetId, err := e.GetColSelection("SUBNET ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), subnetId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e VPCSubnets) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e VPCSubnets) Render() {
	model, err := e.repo.ListSubnets(e.subnetIds)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, subnetId, state, availabilityZone, ipv4Cidr, vpcId string
		if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
			name = n
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
