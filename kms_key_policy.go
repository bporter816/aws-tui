package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type KmsKeyPolicy struct {
	*Text
	kmsClient *kms.Client
	keyId     string
}

func NewKmsKeyPolicy(kmsClient *kms.Client, keyId string) *KmsKeyPolicy {
	k := &KmsKeyPolicy{
		Text:      NewText(true, "json"),
		kmsClient: kmsClient,
		keyId:     keyId,
	}
	return k
}

func (k KmsKeyPolicy) GetName() string {
	return fmt.Sprintf("KMS | %v | Key Policy", k.keyId)
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
