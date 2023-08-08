package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/bporter816/aws-tui/ui"
)

type KmsKeyGrants struct {
	*ui.Table
	kmsClient *kms.Client
	keyId     string
	app       *Application
}

func NewKmsKeyGrants(kmsClient *kms.Client, keyId string, app *Application) *KmsKeyGrants {
	k := &KmsKeyGrants{
		Table: ui.NewTable([]string{
			"NAME",
			"OPERATIONS",
			"GRANTEE PRINCIPAL",
			"RETIRING PRINCIPAL",
		}, 1, 0),
		kmsClient: kmsClient,
		keyId:     keyId,
		app:       app,
	}
	return k
}

func (k KmsKeyGrants) GetService() string {
	return "KMS"
}

func (k KmsKeyGrants) GetLabels() []string {
	return []string{k.keyId, "Grants"}
}

func (k KmsKeyGrants) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsKeyGrants) Render() {
	var grants []kmsTypes.GrantListEntry
	pg := kms.NewListGrantsPaginator(
		k.kmsClient,
		&kms.ListGrantsInput{
			KeyId: aws.String(k.keyId),
		},
	)
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		grants = append(grants, out.Grants...)
	}

	var data [][]string
	for _, v := range grants {
		var name, operations, granteePrincipal, retiringPrincipal string
		if v.Name != nil {
			name = *v.Name
		}
		operations = joinGrantOperations(v.Operations, ", ")
		if v.GranteePrincipal != nil {
			granteePrincipal = *v.GranteePrincipal
		}
		if v.RetiringPrincipal != nil {
			retiringPrincipal = *v.RetiringPrincipal
		}
		data = append(data, []string{
			name,
			operations,
			granteePrincipal,
			retiringPrincipal,
		})
	}
	k.SetData(data)
}
