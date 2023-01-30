package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"strconv"
)

type KmsKeys struct {
	client *kms.Client
}

/*
type KmsKeysEntry struct {
	Id      string
	Aliases string
}
*/

func NewKmsKeys(c *kms.Client) *KmsKeys {
	return &KmsKeys{
		client: c,
	}
}

func (r KmsKeys) GetHeaders() []string {
	return []string{
		"ID",
		"ALIASES",
		"ENABLED",
		"STATE",
		"SPEC",
		"USAGE",
		"REGIONALITY",
	}
}

func (r *KmsKeys) Render() ([][]string, error) {
	aliasMap := make(map[string][]string)
	aliasesPaginator := kms.NewListAliasesPaginator(r.client, &kms.ListAliasesInput{})
	for aliasesPaginator.HasMorePages() {
		out, err := aliasesPaginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, v := range out.Aliases {
			if v.TargetKeyId != nil {
				aliasMap[*v.TargetKeyId] = append(aliasMap[*v.TargetKeyId], *v.AliasName)
			}
		}
	}

	var keys []kmsTypes.KeyListEntry
	keysPaginator := kms.NewListKeysPaginator(r.client, &kms.ListKeysInput{})
	for keysPaginator.HasMorePages() {
		out, err := keysPaginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		keys = append(keys, out.Keys...)
	}

	var ret [][]string
	for _, v := range keys {
		aliases, ok := aliasMap[*v.KeyId]
		var alias string
		if ok {
			alias = aliases[0]
			if len(aliases) > 1 {
				alias += fmt.Sprintf(" + %v more", len(aliases)-1)
			}
		}

		out, err := r.client.DescribeKey(context.TODO(), &kms.DescribeKeyInput{KeyId: v.KeyId})
		if err != nil {
			return nil, err
		}
		var regionality string
		if *out.KeyMetadata.MultiRegion && out.KeyMetadata.MultiRegionConfiguration != nil {
			regionality = string(out.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType)
		}

		ret = append(ret, []string{
			*v.KeyId,
			alias,
			strconv.FormatBool(out.KeyMetadata.Enabled),
			string(out.KeyMetadata.KeyState),
			string(out.KeyMetadata.KeySpec),
			string(out.KeyMetadata.KeyUsage),
			regionality,
		})
	}
	return ret, nil
}
