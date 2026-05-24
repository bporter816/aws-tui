package internal

import (
	"context"
	"fmt"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
)

type IAMPolicy struct {
	*ui.Text
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

func (i IAMPolicy) GetService() string { return "IAM" }

func (i IAMPolicy) GetLabels() []string {
	if i.policyType == model.IAMPolicyTypePermissionsBoundary || i.policyType == model.IAMPolicyTypeAssumeRolePolicy {
		return []string{i.identityName, string(i.policyType)}
	}
	return []string{i.policyName, "Policy Document"}
}

func (i IAMPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMPolicy) Render(ctx context.Context) error {
	var policy string
	var err error
	switch i.policyType {
	case model.IAMPolicyTypeManaged:
		policy, err = i.repo.GetIAMManagedPolicy(ctx, i.policyArn)
	case model.IAMPolicyTypeInline:
		policy, err = i.repo.GetIAMInlinePolicy(ctx, i.identityType, i.identityName, i.policyName)
	case model.IAMPolicyTypePermissionsBoundary:
		policy, err = i.repo.GetIAMPermissionsBoundary(ctx, i.identityName, i.identityType)
	case model.IAMPolicyTypeAssumeRolePolicy:
		policy, err = i.repo.GetIAMAssumeRolePolicy(ctx, i.identityName)
	default:
		return fmt.Errorf("invalid policy type: %v", i.policyType)
	}
	if err != nil {
		return err
	}
	i.SetText(policy)
	return nil
}
