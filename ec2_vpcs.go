package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type EC2VPCs struct {
	*ui.Table
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2VPCs(ec2Client *ec2.Client, app *Application) *EC2VPCs {
	e := &EC2VPCs{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"STATE",
			"IPV4 CIDR",
		}, 1, 0),
		ec2Client: ec2Client,
		app:       app,
	}
	return e
}

func (e EC2VPCs) GetService() string {
	return "EC2"
}

func (e EC2VPCs) GetLabels() []string {
	return []string{"VPCs"}
}

func (e EC2VPCs) tagsHandler() {
	vpcId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2VPCTags(e.ec2Client, vpcId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2VPCs) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2VPCs) Render() {
	pg := ec2.NewDescribeVpcsPaginator(
		e.ec2Client,
		&ec2.DescribeVpcsInput{},
	)
	var vpcs []ec2Types.Vpc
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		vpcs = append(vpcs, out.Vpcs...)
	}

	var data [][]string
	for _, v := range vpcs {
		var name, id, state, ipv4CIDR string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		} else {
			name = "-"
		}
		if v.VpcId != nil {
			id = *v.VpcId
		}
		state = TitleCase(string(v.State))
		if v.CidrBlock != nil {
			ipv4CIDR = *v.CidrBlock
		}
		data = append(data, []string{
			name,
			id,
			state,
			ipv4CIDR,
		})
	}
	e.SetData(data)
}
