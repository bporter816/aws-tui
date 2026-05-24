package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/bporter816/aws-tui/internal/model"
)

type KMS struct {
	kmsClient *kms.Client
}

func NewKMS(kmsClient *kms.Client) *KMS {
	return &KMS{
		kmsClient: kmsClient,
	}
}

func (k KMS) getAliasMap(ctx context.Context) (map[string][]string, error) {
	aliasMap := make(map[string][]string)
	aliasesPg := kms.NewListAliasesPaginator(
		k.kmsClient,
		&kms.ListAliasesInput{},
	)
	for aliasesPg.HasMorePages() {
		out, err := aliasesPg.NextPage(ctx)
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

func (k KMS) describeKey(ctx context.Context, keyId string) (kmsTypes.KeyMetadata, error) {
	out, err := k.kmsClient.DescribeKey(
		ctx,
		&kms.DescribeKeyInput{
			KeyId: aws.String(keyId),
		},
	)
	if err != nil {
		return kmsTypes.KeyMetadata{}, err
	}
	return *out.KeyMetadata, nil
}

func (k KMS) ListKeys(ctx context.Context) ([]model.KMSKey, error) {
	aliasMap, err := k.getAliasMap(ctx)
	if err != nil {
		return []model.KMSKey{}, err
	}
	keysPg := kms.NewListKeysPaginator(
		k.kmsClient,
		&kms.ListKeysInput{},
	)
	var keys []model.KMSKey
	for keysPg.HasMorePages() {
		out, err := keysPg.NextPage(ctx)
		if err != nil {
			return []model.KMSKey{}, err
		}
		for _, v := range out.Keys {
			if v.KeyId != nil {
				meta, err := k.describeKey(ctx, *v.KeyId)
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

func (k KMS) ListGrants(ctx context.Context, keyId string) ([]model.KMSGrant, error) {
	pg := kms.NewListGrantsPaginator(
		k.kmsClient,
		&kms.ListGrantsInput{
			KeyId: aws.String(keyId),
		},
	)
	var grants []model.KMSGrant
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []model.KMSGrant{}, err
		}
		for _, v := range out.Grants {
			grants = append(grants, model.KMSGrant(v))
		}
	}
	return grants, nil
}

func (k KMS) GetKeyPolicy(ctx context.Context, keyId string) (string, error) {
	out, err := k.kmsClient.GetKeyPolicy(
		ctx,
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

func (k KMS) ListTags(ctx context.Context, resourceId string) (model.Tags, error) {
	pg := kms.NewListResourceTagsPaginator(
		k.kmsClient,
		&kms.ListResourceTagsInput{
			KeyId: aws.String(resourceId),
		},
	)
	var tags model.Tags
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return model.Tags{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.TagKey, Value: *v.TagValue})
		}
	}
	return tags, nil
}

func (k KMS) ListCustomKeyStores(ctx context.Context) ([]model.KMSCustomKeyStore, error) {
	pg := kms.NewDescribeCustomKeyStoresPaginator(
		k.kmsClient,
		&kms.DescribeCustomKeyStoresInput{},
	)
	var keyStores []model.KMSCustomKeyStore
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []model.KMSCustomKeyStore{}, err
		}
		for _, v := range out.CustomKeyStores {
			keyStores = append(keyStores, model.KMSCustomKeyStore(v))
		}
	}
	return keyStores, nil
}
