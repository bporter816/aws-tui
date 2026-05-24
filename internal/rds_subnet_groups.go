package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewRDSSubnetGroups(repo *repo.RDS, ec2Repo *repo.EC2, app *Application) *TableView[model.RDSSubnetGroup] {
	return NewTableView(
		"RDS",
		[]string{"Subnet Groups"},
		[]string{"NAME", "STATUS", "SUBNETS", "VPC ID", "NETWORK TYPES", "DESCRIPTION"},
		repo.ListSubnetGroups,
		func(v model.RDSSubnetGroup) []string {
			var networkTypes string
			for i, n := range v.SupportedNetworkTypes {
				networkTypes += utils.AutoCase(n)
				if i < len(v.SupportedNetworkTypes)-1 {
					networkTypes += ", "
				}
			}
			return []string{
				utils.DerefString(v.DBSubnetGroupName, ""),
				utils.DerefString(v.SubnetGroupStatus, ""),
				strconv.Itoa(len(v.Subnets)),
				utils.DerefString(v.VpcId, ""),
				networkTypes,
				utils.DerefString(v.DBSubnetGroupDescription, ""),
			}
		},
		func(v *TableView[model.RDSSubnetGroup]) []KeyAction {
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
						if name := item.DBSubnetGroupName; name != nil {
							app.AddAndSwitch(NewVPCSubnets(ec2Repo, subnetIds, *name, app))
						}
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.DBSubnetGroupArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.DBSubnetGroupArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("RDS.Subnet Groups", func(r Repos, app *Application) Component {
		return NewRDSSubnetGroups(r.RDS, r.EC2, app)
	})
}
