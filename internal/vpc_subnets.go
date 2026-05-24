package internal

import (
	"fmt"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewVPCSubnets(repo *repo.EC2, subnetIds []string, label string, app *Application) *TableView[model.EC2Subnet] {
	var labels []string
	if len(subnetIds) > 0 {
		labels = []string{label, "Subnets"}
	} else {
		labels = []string{"Subnets"}
	}
	return NewTableView(
		"VPC",
		labels,
		[]string{"NAME", "SUBNET ID", "STATE", "AVAILABILITY ZONE", "IPV4 CIDR", "VPC ID"},
		func() ([]model.EC2Subnet, error) {
			return repo.ListSubnets(subnetIds)
		},
		func(v model.EC2Subnet) []string {
			var name, availabilityZone string
			if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
				name = n
			}
			if v.AvailabilityZone != nil {
				availabilityZone = *v.AvailabilityZone
				if v.AvailabilityZoneId != nil {
					availabilityZone += fmt.Sprintf(" (%v)", *v.AvailabilityZoneId)
				}
			}
			return []string{
				name,
				utils.DerefString(v.SubnetId, ""),
				utils.TitleCase(string(v.State)),
				availabilityZone,
				utils.DerefString(v.CidrBlock, ""),
				utils.DerefString(v.VpcId, ""),
			}
		},
		func(v *TableView[model.EC2Subnet]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						subnetId, err := v.GetColSelection("SUBNET ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), subnetId, app))
					},
				},
			}
		},
	)
}

func init() {
	register("VPC.Subnets", func(r Repos, app *Application) Component {
		return NewVPCSubnets(r.EC2, []string{}, "", app)
	})
}
