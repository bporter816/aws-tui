package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/ui"
)

type EC2KeyPairPubKey struct {
	*ui.Text
	ec2Client *ec2.Client
	keyPairId string
	app       *Application
}

func NewEC2KeyPairPubKey(ec2Client *ec2.Client, keyPairId string, app *Application) *EC2KeyPairPubKey {
	e := &EC2KeyPairPubKey{
		Text:      ui.NewText(false, ""),
		ec2Client: ec2Client,
		keyPairId: keyPairId,
		app:       app,
	}
	return e
}

func (e EC2KeyPairPubKey) GetService() string {
	return "EC2"
}

func (e EC2KeyPairPubKey) GetLabels() []string {
	return []string{e.keyPairId, "Public Key"}
}

func (e EC2KeyPairPubKey) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2KeyPairPubKey) Render() {
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

	fmt.Printf("got %v pairs\n", len(out.KeyPairs))
	if len(out.KeyPairs) != 1 {
		panic("should get exactly one key pair")
	}

	var pubKey string
	if out.KeyPairs[0].PublicKey != nil {
		pubKey = *out.KeyPairs[0].PublicKey
	}
	e.SetText(pubKey)
}
