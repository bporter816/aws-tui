package internal

import (
	"fmt"

	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewKmsCustomKeyStores(repo *repo.KMS, app *Application) *TableView[model.KMSCustomKeyStore] {
	return NewTableView(
		"KMS",
		[]string{"Custom Key Stores"},
		[]string{"NAME", "TYPE", "CONNECTION STATUS", "CREATED"},
		repo.ListCustomKeyStores,
		func(v model.KMSCustomKeyStore) []string {
			var created, connection string
			if v.ConnectionState == kmsTypes.ConnectionStateTypeFailed {
				connection = fmt.Sprintf("%v (%v)", v.ConnectionState, v.ConnectionErrorCode)
			} else {
				connection = string(v.ConnectionState)
			}
			if v.CreationDate != nil {
				created = v.CreationDate.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.CustomKeyStoreName, ""),
				string(v.CustomKeyStoreType),
				connection,
				created,
			}
		},
		nil,
	)
}
