package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type KmsKeyGrants struct {
	*ui.Table
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
	model, err := k.repo.ListGrants(k.keyId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
