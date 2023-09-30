package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ElasticacheSubnetGroups struct {
	*ui.Table
	repo    *repo.Elasticache
	ec2Repo *repo.EC2
	app     *Application
	model   []model.ElasticacheSubnetGroup
}

func NewElasticacheSubnetGroups(repo *repo.Elasticache, ec2Repo *repo.EC2, app *Application) *ElasticacheSubnetGroups {
	e := &ElasticacheSubnetGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"SUBNETS",
			"VPC ID",
			"NETWORK TYPES",
			"DESCRIPTION",
		}, 1, 0),
		repo:    repo,
		ec2Repo: ec2Repo,
		app:     app,
	}
	return e
}

func (e ElasticacheSubnetGroups) GetService() string {
	return "Elasticache"
}

func (e ElasticacheSubnetGroups) GetLabels() []string {
	return []string{"Subnet Groups"}
}

func (e ElasticacheSubnetGroups) subnetsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	var subnetIds []string
	for _, v := range e.model[row-1].Subnets {
		if v.SubnetIdentifier != nil {
			subnetIds = append(subnetIds, *v.SubnetIdentifier)
		}
	}
	if e.model[row-1].CacheSubnetGroupName == nil {
		return
	}
	subnetsView := NewEC2Subnets(e.ec2Repo, subnetIds, *e.model[row-1].CacheSubnetGroupName, e.app)
	e.app.AddAndSwitch(subnetsView)
}

func (e ElasticacheSubnetGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if e.model[row-1].ARN == nil || e.model[row-1].CacheSubnetGroupName == nil {
		return
	}
	tagsView := NewElasticacheTags(e.repo, ElasticacheResourceTypeSubnetGroup, *e.model[row-1].ARN, *e.model[row-1].CacheSubnetGroupName, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ElasticacheSubnetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			Description: "Subnets",
			Action:      e.subnetsHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElasticacheSubnetGroups) Render() {
	model, err := e.repo.ListSubnetGroups()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var name, subnets, vpcId, networkTypes, description string
		if v.CacheSubnetGroupName != nil {
			name = *v.CacheSubnetGroupName
		}
		subnets = strconv.Itoa(len(v.Subnets))
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		for i, n := range v.SupportedNetworkTypes {
			networkTypes += string(n)
			if i < len(v.SupportedNetworkTypes)-1 {
				networkTypes += ", "
			}
		}
		if v.CacheSubnetGroupDescription != nil {
			description = *v.CacheSubnetGroupDescription
		}
		data = append(data, []string{
			name,
			subnets,
			vpcId,
			networkTypes,
			description,
		})
	}
	e.SetData(data)
}
