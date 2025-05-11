package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type IAMPolicy struct {
	*ui.Text
	view.IAM
	repo         *repo.IAM
	identityType model.IAMIdentityType
	policyType   model.IAMPolicyType
	identityName string
	policyName   string
	policyArn    string // only for managed policies
	app          *Application
}

func NewIAMPolicy(
	repo *repo.IAM,
	identityType model.IAMIdentityType,
	policyType model.IAMPolicyType,
	identityName string,
	policyName string,
	policyArn string,
	app *Application,
) *IAMPolicy {
	i := &IAMPolicy{
		Text:         ui.NewText(true, "json"),
		repo:         repo,
		identityType: identityType,
		policyType:   policyType,
		identityName: identityName,
		policyName:   policyName,
		policyArn:    policyArn,
		app:          app,
	}
	return i
}

func (i IAMPolicy) GetLabels() []string {
	if i.policyType == model.IAMPolicyTypePermissionsBoundary || i.policyType == model.IAMPolicyTypeAssumeRolePolicy {
		return []string{i.identityName, string(i.policyType)}
	} else if i.policyType == model.IAMPolicyTypeManaged || i.policyType == model.IAMPolicyTypeInline {
		return []string{i.policyName, "Policy Document"}
	} else {
		panic("invalid policy type")
	}
}

func (i IAMPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMPolicy) Render() {
	var policy string
	var err error
	switch i.policyType {
	case model.IAMPolicyTypeManaged:
		policy, err = i.repo.GetIAMManagedPolicy(i.policyArn)
	case model.IAMPolicyTypeInline:
		policy, err = i.repo.GetIAMInlinePolicy(i.identityType, i.identityName, i.policyName)
	case model.IAMPolicyTypePermissionsBoundary:
		policy, err = i.repo.GetIAMPermissionsBoundary(i.identityName, i.identityType)
	case model.IAMPolicyTypeAssumeRolePolicy:
		policy, err = i.repo.GetIAMAssumeRolePolicy(i.identityName)
	}
	if err != nil {
		panic(err)
	}
	i.SetText(policy)
}
