package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type IAMPolicy struct {
	*ui.Text
	repo         *repo.IAM
	identityType model.IAMIdentityType
	policyType   IAMPolicyType
	identityName string
	policyName   string
	policyArn    string // only for managed policies
	app          *Application
}

func NewIAMPolicy(
	repo *repo.IAM,
	identityType model.IAMIdentityType,
	policyType IAMPolicyType,
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

func (i IAMPolicy) GetService() string {
	return "IAM"
}

func (i IAMPolicy) GetLabels() []string {
	if i.policyType == IAMPolicyTypePermissionsBoundary || i.policyType == IAMPolicyTypeAssumeRolePolicy {
		return []string{i.identityName, string(i.policyType)}
	} else if i.policyType == IAMPolicyTypeManaged || i.policyType == IAMPolicyTypeInline {
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
	case IAMPolicyTypeManaged:
		policy, err = i.repo.GetIAMManagedPolicy(i.policyArn)
	case IAMPolicyTypeInline:
		policy, err = i.repo.GetIAMInlinePolicy(i.identityType, i.identityName, i.policyName)
	case IAMPolicyTypePermissionsBoundary:
		policy, err = i.repo.GetIAMPermissionsBoundary(i.identityName, i.identityType)
	case IAMPolicyTypeAssumeRolePolicy:
		policy, err = i.repo.GetIAMAssumeRolePolicy(i.identityName)
	}
	if err != nil {
		panic(err)
	}
	i.SetText(policy)
}
