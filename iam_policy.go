package main

import (
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/bporter816/aws-tui/ui"
)

type IAMPolicy struct {
	*ui.Text
	iamClient    *iam.Client
	app          *Application
	identityType IAMIdentityType
	policyType   IAMPolicyType
	identityName string
	policyName   string
	policyArn    string // only for managed policies
}

func NewIAMPolicy(
	iamClient *iam.Client,
	app *Application,
	identityType IAMIdentityType,
	policyType IAMPolicyType,
	identityName string,
	policyName string,
	policyArn string,
) *IAMPolicy {
	i := &IAMPolicy{
		Text:         ui.NewText(true, "json"),
		iamClient:    iamClient,
		app:          app,
		identityType: identityType,
		policyType:   policyType,
		identityName: identityName,
		policyName:   policyName,
		policyArn:    policyArn,
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
		policy, err = getIAMManagedPolicy(i.iamClient, i.policyArn)
	case IAMPolicyTypeInline:
		policy, err = getIAMInlinePolicy(i.iamClient, i.identityType, i.identityName, i.policyName)
	case IAMPolicyTypePermissionsBoundary:
		policy, err = getIAMPermissionsBoundary(i.iamClient, i.identityName, i.identityType)
	case IAMPolicyTypeAssumeRolePolicy:
		policy, err = getIAMAssumeRolePolicy(i.iamClient, i.identityName)
	}
	if err != nil {
		panic(err)
	}
	i.SetText(policy)
}
