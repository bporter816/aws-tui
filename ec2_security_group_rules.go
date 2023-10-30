package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type EC2SecurityGroupRules struct {
	*ui.Table
	repo *repo.EC2
	sgId string
	app  *Application
}

func NewEC2SecurityGroupRules(repo *repo.EC2, sgId string, app *Application) *EC2SecurityGroupRules {
	e := &EC2SecurityGroupRules{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"TYPE",
			"PROTOCOL",
			"PORTS",
			"SRC/DST",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		sgId: sgId,
		app:  app,
	}
	return e
}

func (e EC2SecurityGroupRules) GetService() string {
	return "EC2"
}

func (e EC2SecurityGroupRules) GetLabels() []string {
	return []string{e.sgId, "Rules"}
}

func (e EC2SecurityGroupRules) tagsHandler() {
	ruleId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SecurityGroupRuleTags(e.repo, ruleId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroupRules) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2SecurityGroupRules) Render() {
	model, err := e.repo.ListSecurityGroupRules(e.sgId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		name := "-"
		var id, ruleType, protocol, ports, cidr, description string
		if ruleName, ok := lookupTag(v.Tags, "Name"); ok {
			name = ruleName
		}
		if v.SecurityGroupRuleId != nil {
			id = *v.SecurityGroupRuleId
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
		if v.Description != nil {
			description = *v.Description
		}
		data = append(data, []string{
			name,
			id,
			ruleType,
			protocol,
			ports,
			cidr,
			description,
		})
	}
	e.SetData(data)
}
