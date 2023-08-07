package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bporter816/aws-tui/ui"
)

type EC2KeyPairs struct {
	*ui.Table
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2KeyPairs(ec2Client *ec2.Client, app *Application) *EC2KeyPairs {
	e := &EC2KeyPairs{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"FINGERPRINT",
			"CREATED",
			"ID",
		}, 1, 0),
		ec2Client: ec2Client,
		app:       app,
	}
	return e
}

func (e EC2KeyPairs) GetService() string {
	return "EC2"
}

func (e EC2KeyPairs) GetLabels() []string {
	return []string{"Key Pairs"}
}

func (e EC2KeyPairs) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2KeyPairs) Render() {
	out, err := e.ec2Client.DescribeKeyPairs(
		context.TODO(),
		&ec2.DescribeKeyPairsInput{},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range out.KeyPairs {
		var name, keyType, fingerprint, created, id string
		if v.KeyName != nil {
			name = *v.KeyName
		}
		keyType = string(v.KeyType)
		if v.KeyFingerprint != nil {
			fingerprint = *v.KeyFingerprint
		}
		if v.CreateTime != nil {
			created = v.CreateTime.Format(DefaultTimeFormat)
		}
		if v.KeyPairId != nil {
			id = *v.KeyPairId
		}
		data = append(data, []string{
			name,
			keyType,
			fingerprint,
			created,
			id,
		})
	}
	e.SetData(data)
}
