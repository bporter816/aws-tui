package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/bporter816/aws-tui/ui"
)

type KmsKeyPolicy struct {
	*ui.Text
	kmsClient *kms.Client
	keyId     string
}

func NewKmsKeyPolicy(kmsClient *kms.Client, keyId string) *KmsKeyPolicy {
	k := &KmsKeyPolicy{
		Text:      ui.NewText(true, "json"),
		kmsClient: kmsClient,
		keyId:     keyId,
	}
	return k
}

func (k KmsKeyPolicy) GetService() string {
	return "KMS"
}

func (k KmsKeyPolicy) GetLabels() []string {
	return []string{k.keyId, "Key Policy"}
}

func (k KmsKeyPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsKeyPolicy) Render() {
	out, err := k.kmsClient.GetKeyPolicy(
		context.TODO(),
		&kms.GetKeyPolicyInput{
			KeyId:      aws.String(k.keyId),
			PolicyName: aws.String("default"),
		},
	)
	if err != nil {
		panic(err)
	}
	var policy string
	if out.Policy != nil {
		policy = *out.Policy
	}
	k.SetText(policy)
}
