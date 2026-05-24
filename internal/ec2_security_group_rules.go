package internal

import (
	"fmt"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEC2SecurityGroupRules(repo *repo.EC2, sgId string, app *Application) *TableView[model.EC2SecurityGroupRule] {
	return NewTableView(
		"EC2",
		[]string{sgId, "Rules"},
		[]string{"NAME", "ID", "TYPE", "PROTOCOL", "PORTS", "SRC/DST", "DESCRIPTION"},
		func() ([]model.EC2SecurityGroupRule, error) {
			return repo.ListSecurityGroupRules(sgId)
		},
		func(v model.EC2SecurityGroupRule) []string {
			name := "-"
			var ruleType, protocol, ports, cidr string
			if ruleName, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
				name = ruleName
			}
			if v.IsEgress != nil {
				if *v.IsEgress {
					ruleType = "Egress"
				} else {
					ruleType = "Ingress"
				}
			}
			if v.IpProtocol != nil {
				proto := *v.IpProtocol
				if proto == "-1" {
					protocol = "All"
				} else {
					protocol = utils.UpperCase(*v.IpProtocol)
				}
			}
			if v.CidrIpv4 != nil {
				cidr = *v.CidrIpv4
			} else if v.CidrIpv6 != nil {
				cidr = *v.CidrIpv6
			} else if v.ReferencedGroupInfo != nil && v.ReferencedGroupInfo.GroupId != nil {
				cidr = *v.ReferencedGroupInfo.GroupId
			}
			if v.FromPort != nil && v.ToPort != nil {
				from, to := *v.FromPort, *v.ToPort
				if from == -1 {
					ports = "All"
				} else if from == to {
					ports = fmt.Sprintf("%v", from)
				} else {
					// TODO handle ICMP case better
					ports = fmt.Sprintf("%v-%v", from, to)
				}
			}
			return []string{
				name,
				utils.DerefString(v.SecurityGroupRuleId, ""),
				ruleType,
				protocol,
				ports,
				cidr,
				utils.DerefString(v.Description, ""),
			}
		},
		func(v *TableView[model.EC2SecurityGroupRule]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						ruleId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), ruleId, app))
					},
				},
			}
		},
	)
}
