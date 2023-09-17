package main

import (
	"fmt"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type KmsCustomKeyStores struct {
	*ui.Table
	repo *repo.KMS
	app  *Application
}

func NewKmsCustomKeyStores(repo *repo.KMS, app *Application) *KmsCustomKeyStores {
	k := &KmsCustomKeyStores{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"CONNECTION STATUS",
			"CREATED",
		}, 1, 0),
		repo: repo,
		app:  app,
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
	model, err := k.repo.ListCustomKeyStores()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
