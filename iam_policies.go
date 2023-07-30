package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type IAMPolicies struct {
	*ui.Table
	iamClient         *iam.Client
	app               *Application
	identityType      IAMIdentityType
	id                string
	managedArns       []string
	numInlinePolicies int
}

func NewIAMPolicies(iamClient *iam.Client, app *Application, identityType IAMIdentityType, id string) *IAMPolicies {
	i := &IAMPolicies{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
		}, 1, 0),
		iamClient:    iamClient,
		app:          app,
		identityType: identityType,
		id:           id,
	}
	return i
}

func (i IAMPolicies) GetService() string {
	return "IAM"
}

func (i IAMPolicies) GetLabels() []string {
	return []string{i.id, "Policies"}
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
	enumVal := IAMPolicyType(policyType)
	var policyDocumentView *IAMPolicy
	if enumVal == IAMPolicyTypeManaged {
		policyDocumentView = NewIAMPolicy(i.iamClient, i.app, i.identityType, enumVal, i.id, policyName, i.managedArns[row-1-i.numInlinePolicies])
	} else {
		policyDocumentView = NewIAMPolicy(i.iamClient, i.app, i.identityType, enumVal, i.id, policyName, "")
	}
	i.app.AddAndSwitch(policyDocumentView)
}

func (i IAMPolicies) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone),
			Description: "Policy Document",
			Action:      i.policyDocumentHandler,
		},
	}
}

func (i *IAMPolicies) Render() {
	var data [][]string

	if i.id != "" {
		// inline policies
		var inlinePolicyNames []string
		switch i.identityType {
		case IAMIdentityTypeUser:
			pg := iam.NewListUserPoliciesPaginator(
				i.iamClient,
				&iam.ListUserPoliciesInput{
					UserName: aws.String(i.id),
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(context.TODO())
				if err != nil {
					panic(err)
				}
				inlinePolicyNames = append(inlinePolicyNames, out.PolicyNames...)
			}
		case IAMIdentityTypeRole:
			pg := iam.NewListRolePoliciesPaginator(
				i.iamClient,
				&iam.ListRolePoliciesInput{
					RoleName: aws.String(i.id),
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(context.TODO())
				if err != nil {
					panic(err)
				}
				inlinePolicyNames = append(inlinePolicyNames, out.PolicyNames...)
			}
		case IAMIdentityTypeGroup:
			pg := iam.NewListGroupPoliciesPaginator(
				i.iamClient,
				&iam.ListGroupPoliciesInput{
					GroupName: aws.String(i.id),
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(context.TODO())
				if err != nil {
					panic(err)
				}
				inlinePolicyNames = append(inlinePolicyNames, out.PolicyNames...)
			}
		default:
			panic("invalid identity type for policy list")
		}

		i.numInlinePolicies = len(inlinePolicyNames)
		for _, v := range inlinePolicyNames {
			data = append(data, []string{
				v,
				string(IAMPolicyTypeInline),
			})
		}
	}

	// managed policies
	if i.id == "" {
		// all policies in account
		var policies []iamTypes.Policy
		pg := iam.NewListPoliciesPaginator(
			i.iamClient,
			&iam.ListPoliciesInput{},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(context.TODO())
			if err != nil {
				panic(err)
			}
			policies = append(policies, out.Policies...)
		}

		i.managedArns = make([]string, len(policies))
		for idx, v := range policies {
			i.managedArns[idx] = *v.Arn
			var name string
			if v.PolicyName != nil {
				name = *v.PolicyName
			}
			data = append(data, []string{
				name,
				string(IAMPolicyTypeManaged),
			})
		}
	} else {
		// policies attached to a user, role, or group
		var attachedPolicies []iamTypes.AttachedPolicy
		switch i.identityType {
		case IAMIdentityTypeUser:
			pg := iam.NewListAttachedUserPoliciesPaginator(
				i.iamClient,
				&iam.ListAttachedUserPoliciesInput{
					UserName: aws.String(i.id),
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(context.TODO())
				if err != nil {
					panic(err)
				}
				attachedPolicies = append(attachedPolicies, out.AttachedPolicies...)
			}
		case IAMIdentityTypeRole:
			pg := iam.NewListAttachedRolePoliciesPaginator(
				i.iamClient,
				&iam.ListAttachedRolePoliciesInput{
					RoleName: aws.String(i.id),
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(context.TODO())
				if err != nil {
					panic(err)
				}
				attachedPolicies = append(attachedPolicies, out.AttachedPolicies...)
			}
		case IAMIdentityTypeGroup:
			pg := iam.NewListAttachedGroupPoliciesPaginator(
				i.iamClient,
				&iam.ListAttachedGroupPoliciesInput{
					GroupName: aws.String(i.id),
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(context.TODO())
				if err != nil {
					panic(err)
				}
				attachedPolicies = append(attachedPolicies, out.AttachedPolicies...)
			}
		default:
			panic("invalid identity type for policy list")
		}

		i.managedArns = make([]string, len(attachedPolicies))
		for idx, v := range attachedPolicies {
			i.managedArns[idx] = *v.PolicyArn
			var name string
			if v.PolicyName != nil {
				name = *v.PolicyName
			}
			data = append(data, []string{
				name,
				string(IAMPolicyTypeManaged),
			})
		}
	}

	i.SetData(data)
}
