package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheSubnetGroups(repo *repo.ElastiCache, ec2Repo *repo.EC2, app *Application) *TableView[model.ElastiCacheSubnetGroup] {
	return NewTableView(
		"ElastiCache",
		[]string{"Subnet Groups"},
		[]string{"NAME", "SUBNETS", "VPC ID", "NETWORK TYPES", "DESCRIPTION"},
		repo.ListSubnetGroups,
		func(v model.ElastiCacheSubnetGroup) []string {
			var networkTypes string
			for i, n := range v.SupportedNetworkTypes {
				networkTypes += utils.AutoCase(string(n))
				if i < len(v.SupportedNetworkTypes)-1 {
					networkTypes += ", "
				}
			}
			return []string{
				utils.DerefString(v.CacheSubnetGroupName, ""),
				strconv.Itoa(len(v.Subnets)),
				utils.DerefString(v.VpcId, ""),
				networkTypes,
				utils.DerefString(v.CacheSubnetGroupDescription, ""),
			}
		},
		func(v *TableView[model.ElastiCacheSubnetGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
					Description: "Subnets",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						var subnetIds []string
						for _, s := range item.Subnets {
							if s.SubnetIdentifier != nil {
								subnetIds = append(subnetIds, *s.SubnetIdentifier)
							}
						}
						if name := item.CacheSubnetGroupName; name != nil {
							app.AddAndSwitch(NewVPCSubnets(ec2Repo, subnetIds, *name, app))
						}
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ARN == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ARN, app))
					},
				},
			}
		},
	)
}
