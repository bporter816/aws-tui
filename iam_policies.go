package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
	// "github.com/gdamore/tcell/v2"
)

type IAMPolicies struct {
	*ui.Table
	iamClient    *iam.Client
	app          *Application
	identityType IAMIdentityType
	id           string
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

func (i IAMPolicies) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMPolicies) Render() {
	var data [][]string

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

	for _, v := range inlinePolicyNames {
		data = append(data, []string{
			v,
			string(IAMPolicyTypeInline),
		})
	}

	// managed policies
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

	for _, v := range attachedPolicies {
		var name string
		if v.PolicyName != nil {
			name = *v.PolicyName
		}
		data = append(data, []string{
			name,
			string(IAMPolicyTypeManaged),
		})
	}

	i.SetData(data)
}
