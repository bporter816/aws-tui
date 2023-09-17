package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/bporter816/aws-tui/model"
)

type KMS struct {
	kmsClient *kms.Client
}

func NewKMS(kmsClient *kms.Client) *KMS {
	return &KMS{
		kmsClient: kmsClient,
	}
}

func (k KMS) getAliasMap() (map[string][]string, error) {
	aliasMap := make(map[string][]string)
	aliasesPg := kms.NewListAliasesPaginator(
		k.kmsClient,
		&kms.ListAliasesInput{},
	)
	for aliasesPg.HasMorePages() {
		out, err := aliasesPg.NextPage(context.TODO())
		if err != nil {
			return map[string][]string{}, err
		}
		for _, v := range out.Aliases {
			if v.TargetKeyId != nil && v.AliasName != nil {
				aliasMap[*v.TargetKeyId] = append(aliasMap[*v.TargetKeyId], *v.AliasName)
			}
		}
	}
	return aliasMap, nil
}

func (k KMS) describeKey(keyId string) (kmsTypes.KeyMetadata, error) {
	out, err := k.kmsClient.DescribeKey(
		context.TODO(),
		&kms.DescribeKeyInput{
			KeyId: aws.String(keyId),
		},
	)
	if err != nil {
		return kmsTypes.KeyMetadata{}, err
	}
	return *out.KeyMetadata, nil
}

func (k KMS) ListKeys() ([]model.KMSKey, error) {
	aliasMap, err := k.getAliasMap()
	if err != nil {
		return []model.KMSKey{}, err
	}
	keysPg := kms.NewListKeysPaginator(
		k.kmsClient,
		&kms.ListKeysInput{},
	)
	var keys []model.KMSKey
	for keysPg.HasMorePages() {
		out, err := keysPg.NextPage(context.TODO())
		if err != nil {
			return []model.KMSKey{}, err
		}
		for _, v := range out.Keys {
			if v.KeyId != nil {
				meta, err := k.describeKey(*v.KeyId)
				if err != nil {
					// TODO handle error
					continue
				}
				m := model.KMSKey{KeyMetadata: meta}
				if a, ok := aliasMap[*v.KeyId]; ok {
					m.Aliases = a
				}
				keys = append(keys, m)
			}
		}
	}
	return keys, nil
}

func (k KMS) ListGrants(keyId string) ([]model.KMSGrant, error) {
	pg := kms.NewListGrantsPaginator(
		k.kmsClient,
		&kms.ListGrantsInput{
			KeyId: aws.String(keyId),
		},
	)
	var grants []model.KMSGrant
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.KMSGrant{}, err
		}
		for _, v := range out.Grants {
			grants = append(grants, model.KMSGrant(v))
		}
	}
	return grants, nil
}

func (k KMS) GetKeyPolicy(keyId string) (string, error) {
	out, err := k.kmsClient.GetKeyPolicy(
		context.TODO(),
		&kms.GetKeyPolicyInput{
			KeyId:      aws.String(keyId),
			PolicyName: aws.String("default"),
		},
	)
	if err != nil || out.Policy == nil {
		return "", err
	}
	return *out.Policy, nil
}

func (k KMS) ListTags(keyId string) ([]model.Tag, error) {
	pg := kms.NewListResourceTagsPaginator(
		k.kmsClient,
		&kms.ListResourceTagsInput{
			KeyId: aws.String(keyId),
		},
	)
	var tags []model.Tag
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.Tag{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.TagKey, Value: *v.TagValue})
		}
	}
	return tags, nil
}
