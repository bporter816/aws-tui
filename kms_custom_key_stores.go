package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type KmsCustomKeyStores struct {
	*ui.Table
	kmsClient *kms.Client
	app       *Application
}

func NewKmsCustomKeyStores(kmsClient *kms.Client, app *Application) *KmsCustomKeyStores {
	k := &KmsCustomKeyStores{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"CONNECTION STATUS",
			"CREATED",
		}, 1, 0),
		kmsClient: kmsClient,
		app:       app,
	}
	return k
}

func (k KmsCustomKeyStores) GetService() string {
	return "KMS"
}

func (k KmsCustomKeyStores) GetLabels() []string {
	return []string{"Custom Key Stores"}
}

func (k KmsCustomKeyStores) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsCustomKeyStores) Render() {
	var keyStores []kmsTypes.CustomKeyStoresListEntry
	pg := kms.NewDescribeCustomKeyStoresPaginator(
		k.kmsClient,
		&kms.DescribeCustomKeyStoresInput{},
	)
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		keyStores = append(keyStores, out.CustomKeyStores...)
	}

	var data [][]string
	for _, v := range keyStores {
		var name, keyStoreType, connection, created string
		if v.CustomKeyStoreName != nil {
			name = *v.CustomKeyStoreName
		}
		keyStoreType = string(v.CustomKeyStoreType)
		if v.ConnectionState == kmsTypes.ConnectionStateTypeFailed {
			connection = fmt.Sprintf("%v (%v)", v.ConnectionState, v.ConnectionErrorCode)
		} else {
			connection = string(v.ConnectionState)
		}
		if v.CreationDate != nil {
			created = v.CreationDate.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			name,
			keyStoreType,
			connection,
			created,
		})
	}
	k.SetData(data)
}
