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

type RDSSubnetGroups struct {
	*ui.Table
	view.RDS
	repo    *repo.RDS
	ec2Repo *repo.EC2
	app     *Application
	model   []model.RDSSubnetGroup
}

func NewRDSSubnetGroups(repo *repo.RDS, ec2Repo *repo.EC2, app *Application) *RDSSubnetGroups {
	r := &RDSSubnetGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"SUBNETS",
			"VPC ID",
			"NETWORK TYPES",
			"DESCRIPTION",
		}, 1, 0),
		repo:    repo,
		ec2Repo: ec2Repo,
		app:     app,
	}
	return r
}

func (r RDSSubnetGroups) GetLabels() []string {
	return []string{"Subnet Groups"}
}

func (r RDSSubnetGroups) subnetsHandler() {
	row, err := r.GetRowSelection()
	if err != nil {
		return
	}
	var subnetIds []string
	for _, v := range r.model[row-1].Subnets {
		if v.SubnetIdentifier != nil {
			subnetIds = append(subnetIds, *v.SubnetIdentifier)
		}
	}
	if name := r.model[row-1].DBSubnetGroupName; name != nil {
		subnetsView := NewVPCSubnets(r.ec2Repo, subnetIds, *name, r.app)
		r.app.AddAndSwitch(subnetsView)
	}
}

func (r RDSSubnetGroups) tagsHandler() {
	row, err := r.GetRowSelection()
	if err != nil {
		return
	}
	if arn := r.model[row-1].DBSubnetGroupArn; arn != nil {
		tagsView := NewTags(r.repo, r.GetService(), *arn, r.app)
		r.app.AddAndSwitch(tagsView)
	}
}

func (r RDSSubnetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			Description: "Subnets",
			Action:      r.subnetsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r *RDSSubnetGroups) Render() {
	model, err := r.repo.ListSubnetGroups()
	if err != nil {
		panic(err)
	}
	r.model = model

	var data [][]string
	for _, v := range model {
		var name, status, subnets, vpcId, networkTypes, desc string
		if v.DBSubnetGroupName != nil {
			name = *v.DBSubnetGroupName
		}
		if v.SubnetGroupStatus != nil {
			status = *v.SubnetGroupStatus
		}
		subnets = strconv.Itoa(len(v.Subnets))
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		for i, n := range v.SupportedNetworkTypes {
			networkTypes += utils.AutoCase(n)
			if i < len(v.SupportedNetworkTypes)-1 {
				networkTypes += ", "
			}
		}
		if v.DBSubnetGroupDescription != nil {
			desc = *v.DBSubnetGroupDescription
		}
		data = append(data, []string{
			name,
			status,
			subnets,
			vpcId,
			networkTypes,
			desc,
		})
	}
	r.SetData(data)
}
