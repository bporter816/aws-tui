package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheSubnetGroups struct {
	*ui.Table
	view.ElastiCache
	repo    *repo.ElastiCache
	ec2Repo *repo.EC2
	app     *Application
	model   []model.ElastiCacheSubnetGroup
}

func NewElastiCacheSubnetGroups(repo *repo.ElastiCache, ec2Repo *repo.EC2, app *Application) *ElastiCacheSubnetGroups {
	e := &ElastiCacheSubnetGroups{
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

func (e ElastiCacheSubnetGroups) GetLabels() []string {
	return []string{"Subnet Groups"}
}

func (e ElastiCacheSubnetGroups) subnetsHandler() {
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
	if name := e.model[row-1].CacheSubnetGroupName; name != nil {
		subnetsView := NewVPCSubnets(e.ec2Repo, subnetIds, *name, e.app)
		e.app.AddAndSwitch(subnetsView)
	}
}

func (e ElastiCacheSubnetGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].ARN; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheSubnetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			Description: "Subnets",
			Action:      e.subnetsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElastiCacheSubnetGroups) Render() {
	model, err := e.repo.ListSubnetGroups()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var networkTypes string
		for i, n := range v.SupportedNetworkTypes {
			networkTypes += utils.AutoCase(string(n))
			if i < len(v.SupportedNetworkTypes)-1 {
				networkTypes += ", "
			}
		}
		data = append(data, []string{
			utils.DerefString(v.CacheSubnetGroupName, ""),
			strconv.Itoa(len(v.Subnets)),
			utils.DerefString(v.VpcId, ""),
			networkTypes,
			utils.DerefString(v.CacheSubnetGroupDescription, ""),
		})
	}
	e.SetData(data)
}
