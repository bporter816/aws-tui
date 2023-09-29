package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"strconv"
)

type ElasticacheSubnetGroups struct {
	*ui.Table
	repo *repo.Elasticache
	app  *Application
}

func NewElasticacheSubnetGroups(repo *repo.Elasticache, app *Application) *ElasticacheSubnetGroups {
	e := &ElasticacheSubnetGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"SUBNETS",
			"VPC ID",
			"NETWORK TYPES",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElasticacheSubnetGroups) GetService() string {
	return "Elasticache"
}

func (e ElasticacheSubnetGroups) GetLabels() []string {
	return []string{"Events"}
}

func (e ElasticacheSubnetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheSubnetGroups) Render() {
	model, err := e.repo.ListSubnetGroups()
	if err != nil {
		panic(err)
	}

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
