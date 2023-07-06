package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KmsKeyTags struct {
	*Table
	kmsClient *kms.Client
	app       *Application
	keyId     string
}

func NewKmsKeyTags(kmsClient *kms.Client, keyId string, app *Application) *KmsKeyTags {
	k := &KmsKeyTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		kmsClient: kmsClient,
		keyId:     keyId,
		app:       app,
	}
	return k
}

func (k KmsKeyTags) GetName() string {
	return fmt.Sprintf("KMS | %v | Tags", k.keyId)
}

func (k KmsKeyTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsKeyTags) Render() {
	var tags []kmsTypes.Tag
	pg := kms.NewListResourceTagsPaginator(
		k.kmsClient,
		&kms.ListResourceTagsInput{
			KeyId: aws.String(k.keyId),
		},
	)
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		tags = append(tags, out.Tags...)
	}

	var data [][]string
	for _, v := range tags {
		data = append(data, []string{
			*v.TagKey,
			*v.TagValue,
		})
	}
	k.SetData(data)
}
