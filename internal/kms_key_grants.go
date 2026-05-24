package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewKmsKeyGrants(repo *repo.KMS, keyId string, app *Application) *TableView[model.KMSGrant] {
	return NewTableView(
		"KMS",
		[]string{keyId, "Grants"},
		[]string{"NAME", "OPERATIONS", "GRANTEE PRINCIPAL", "RETIRING PRINCIPAL"},
		func() ([]model.KMSGrant, error) {
			return repo.ListGrants(keyId)
		},
		func(v model.KMSGrant) []string {
			return []string{
				utils.DerefString(v.Name, ""),
				utils.JoinKMSGrantOperations(v.Operations, ", "),
				utils.DerefString(v.GranteePrincipal, ""),
				utils.DerefString(v.RetiringPrincipal, ""),
			}
		},
		nil,
	)
}
