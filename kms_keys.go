package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type KmsKeys struct {
	*Table
	kmsClient *kms.Client
	app       *Application
}

func NewKmsKeys(kmsClient *kms.Client, app *Application) *KmsKeys {
	k := &KmsKeys{
		Table: NewTable([]string{
			"ID",
			"ALIASES",
			"DESCRIPTION",
			"ENABLED",
			"STATE",
			"SPEC",
			"USAGE",
			"REGIONALITY",
		}, 1, 0),
		kmsClient: kmsClient,
		app:       app,
	}
	return k
}

func (k KmsKeys) GetName() string {
	return "KMS | Keys"
}

func (k KmsKeys) keyPolicyHandler() {
	keyId, err := k.GetColSelection("ID")
	if err != nil {
		panic(err)
	}
	policyView := NewKmsKeyPolicy(k.kmsClient, keyId)
	k.app.AddAndSwitch("kms.policy", policyView)
}

func (k KmsKeys) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Key Policy",
			Action:      k.keyPolicyHandler,
		},
	}
}

func (k KmsKeys) Render() {
	aliasMap := make(map[string][]string)
	aliasesPg := kms.NewListAliasesPaginator(
		k.kmsClient,
		&kms.ListAliasesInput{},
	)
	for aliasesPg.HasMorePages() {
		out, err := aliasesPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, v := range out.Aliases {
			if v.TargetKeyId != nil {
				aliasMap[*v.TargetKeyId] = append(aliasMap[*v.TargetKeyId], *v.AliasName)
			}
		}
	}

	var keys []kmsTypes.KeyListEntry
	keysPg := kms.NewListKeysPaginator(
		k.kmsClient,
		&kms.ListKeysInput{},
	)
	for keysPg.HasMorePages() {
		out, err := keysPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		keys = append(keys, out.Keys...)
	}

	var data [][]string
	for _, v := range keys {
		aliases, ok := aliasMap[*v.KeyId]
		var alias string
		if ok {
			alias = aliases[0]
			if len(aliases) > 1 {
				alias += fmt.Sprintf(" + %v more", len(aliases)-1)
			}
		}

		out, err := k.kmsClient.DescribeKey(
			context.TODO(),
			&kms.DescribeKeyInput{
				KeyId: v.KeyId,
			},
		)
		if err != nil {
			// panic(err)
			data = append(data, []string{
				*v.KeyId,
				"Unauthorized",
				"-",
				"-",
				"-",
				"-",
				"-",
				"-",
			})
			continue
		}
		var regionality string
		if *out.KeyMetadata.MultiRegion && out.KeyMetadata.MultiRegionConfiguration != nil {
			regionality = string(out.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType)
		}

		data = append(data, []string{
			*v.KeyId,
			alias,
			*out.KeyMetadata.Description,
			strconv.FormatBool(out.KeyMetadata.Enabled),
			string(out.KeyMetadata.KeyState),
			string(out.KeyMetadata.KeySpec),
			string(out.KeyMetadata.KeyUsage),
			regionality,
		})
	}
	k.SetData(data)
}
