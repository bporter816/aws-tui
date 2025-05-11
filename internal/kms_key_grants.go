package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type KmsKeyGrants struct {
	*ui.Table
	view.KMS
	repo  *repo.KMS
	keyId string
	app   *Application
}

func NewKmsKeyGrants(repo *repo.KMS, keyId string, app *Application) *KmsKeyGrants {
	k := &KmsKeyGrants{
		Table: ui.NewTable([]string{
			"NAME",
			"OPERATIONS",
			"GRANTEE PRINCIPAL",
			"RETIRING PRINCIPAL",
		}, 1, 0),
		repo:  repo,
		keyId: keyId,
		app:   app,
	}
	return k
}

func (k KmsKeyGrants) GetLabels() []string {
	return []string{k.keyId, "Grants"}
}

func (k KmsKeyGrants) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsKeyGrants) Render() {
	model, err := k.repo.ListGrants(k.keyId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			utils.JoinKMSGrantOperations(v.Operations, ", "),
			utils.DerefString(v.GranteePrincipal, ""),
			utils.DerefString(v.RetiringPrincipal, ""),
		})
	}
	k.SetData(data)
}
