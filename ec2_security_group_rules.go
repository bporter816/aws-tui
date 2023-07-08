package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EC2SecurityGroupRules struct {
	*Table
	ec2Client *ec2.Client
	sgId      string
	app       *Application
}

func NewEC2SecurityGroupRules(ec2Client *ec2.Client, sgId string, app *Application) *EC2SecurityGroupRules {
	e := &EC2SecurityGroupRules{
		Table: NewTable([]string{
			"NAME",
			"ID",
			"TYPE",
			"PROTOCOL",
			"PORTS",
			"SRC/DST",
			"DESCRIPTION",
		}, 1, 0),
		ec2Client: ec2Client,
		sgId:      sgId,
		app:       app,
	}
	return e
}

func (e EC2SecurityGroupRules) GetService() string {
	return "EC2"
}

func (e EC2SecurityGroupRules) GetLabels() []string {
	return []string{"Security Groups", e.sgId, "Rules"}
}

func (e EC2SecurityGroupRules) tagsHandler() {
	ruleId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SecurityGroupRuleTags(e.ec2Client, e.sgId, ruleId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroupRules) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2SecurityGroupRules) Render() {
	pg := ec2.NewDescribeSecurityGroupRulesPaginator(
		e.ec2Client,
		&ec2.DescribeSecurityGroupRulesInput{
			Filters: []ec2Types.Filter{
				ec2Types.Filter{
					Name:   aws.String("group-id"),
					Values: []string{e.sgId},
				},
			},
		},
	)
	var sgRules []ec2Types.SecurityGroupRule
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		sgRules = append(sgRules, out.SecurityGroupRules...)
	}

	caser := cases.Upper(language.English)
	var data [][]string
	for _, v := range sgRules {
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
				protocol = caser.String(*v.IpProtocol)
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
