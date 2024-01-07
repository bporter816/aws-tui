package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type IAMPolicies struct {
	*ui.Table
	view.IAM
	repo         *repo.IAM
	identityType model.IAMIdentityType
	id           *string
	app          *Application
	model        []model.IAMPolicy
}

func NewIAMPolicies(repo *repo.IAM, identityType model.IAMIdentityType, id *string, app *Application) *IAMPolicies {
	i := &IAMPolicies{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
		}, 1, 0),
		repo:         repo,
		identityType: identityType,
		id:           id,
		app:          app,
	}
	return i
}

func (i IAMPolicies) GetLabels() []string {
	return []string{*i.id, "Policies"}
}

func (i IAMPolicies) policyDocumentHandler() {
	row, err := i.GetRowSelection()
	if err != nil {
		return
	}
	policyName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	policyType, err := i.GetColSelection("TYPE")
	if err != nil {
		return
	}
	enumVal := model.IAMPolicyType(policyType)
	var policyDocumentView *IAMPolicy
	if enumVal == model.IAMPolicyTypeManaged {
		if i.model[row-1].Arn == nil {
			return
		}
		policyDocumentView = NewIAMPolicy(i.repo, i.identityType, enumVal, *i.id, policyName, *i.model[row-1].Arn, i.app)
	} else {
		policyDocumentView = NewIAMPolicy(i.repo, i.identityType, enumVal, *i.id, policyName, "", i.app)
	}
	i.app.AddAndSwitch(policyDocumentView)
}

func (i IAMPolicies) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone),
			Description: "Policy Document",
			Action:      i.policyDocumentHandler,
		},
	}
}

func (i *IAMPolicies) Render() {
	model, err := i.repo.ListPolicies(i.id, i.identityType)
	if err != nil {
		panic(err)
	}
	i.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Name,
			string(v.PolicyType),
		})
	}
	i.SetData(data)
}
