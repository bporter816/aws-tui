package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EC2VPCs struct {
	*Table
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2VPCs(ec2Client *ec2.Client, app *Application) *EC2VPCs {
	e := &EC2VPCs{
		Table: NewTable([]string{
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

func (e EC2VPCs) GetName() string {
	return "EC2 | VPCs"
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

	caser := cases.Title(language.English)
	var data [][]string
	for _, v := range vpcs {
		name := "-"
		var id, state, ipv4CIDR string
		if v.VpcId != nil {
			id = *v.VpcId
			if vpcName, ok := lookupTag(v.Tags, "Name"); ok {
				name = vpcName
			}
		}
		state = caser.String(string(v.State))
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

func lookupTag(tags []ec2Types.Tag, key string) (string, bool) {
	for _, v := range tags {
		if *v.Key == key {
			return *v.Value, true
		}
	}
	return "", false
}
