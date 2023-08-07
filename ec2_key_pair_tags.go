package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/ui"
)

type EC2KeyPairTags struct {
	*ui.Table
	ec2Client *ec2.Client
	keyPairId string
	app       *Application
}

func NewEC2KeyPairTags(ec2Client *ec2.Client, keyPairId string, app *Application) *EC2KeyPairTags {
	e := &EC2KeyPairTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client: ec2Client,
		keyPairId: keyPairId,
		app:       app,
	}
	return e
}

func (e EC2KeyPairTags) GetService() string {
	return "EC2"
}

func (e EC2KeyPairTags) GetLabels() []string {
	return []string{e.keyPairId, "Tags"}
}

func (e EC2KeyPairTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2KeyPairTags) Render() {
	out, err := e.ec2Client.DescribeKeyPairs(
		context.TODO(),
		&ec2.DescribeKeyPairsInput{
			IncludePublicKey: aws.Bool(true),
			Filters: []ec2Types.Filter{
				ec2Types.Filter{
					Name:   aws.String("key-pair-id"),
					Values: []string{e.keyPairId},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}

	if len(out.KeyPairs) != 1 {
		panic("should get exactly one key pair")
	}

	var data [][]string
	for _, v := range out.KeyPairs[0].Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
