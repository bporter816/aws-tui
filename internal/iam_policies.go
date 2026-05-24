package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/gdamore/tcell/v2"
)

func NewIAMPolicies(repo *repo.IAM, identityType model.IAMIdentityType, id *string, app *Application) *TableView[model.IAMPolicy] {
	return NewTableView(
		"IAM",
		[]string{*id, "Policies"},
		[]string{"NAME", "TYPE"},
		func(ctx context.Context) ([]model.IAMPolicy, error) {
			return repo.ListPolicies(ctx, id, identityType)
		},
		func(v model.IAMPolicy) []string {
			return []string{
				v.Name,
				string(v.PolicyType),
			}
		},
		func(v *TableView[model.IAMPolicy]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone),
					Description: "Policy Document",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						policyName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						policyType, err := v.GetColSelection("TYPE")
						if err != nil {
							return
						}
						enumVal := model.IAMPolicyType(policyType)
						var policyDocumentView *IAMPolicy
						if enumVal == model.IAMPolicyTypeManaged {
							if item.Arn == nil {
								return
							}
							policyDocumentView = NewIAMPolicy(repo, identityType, enumVal, *id, policyName, *item.Arn, app)
						} else {
							policyDocumentView = NewIAMPolicy(repo, identityType, enumVal, *id, policyName, "", app)
						}
						app.AddAndSwitch(policyDocumentView)
					},
				},
			}
		},
	)
}

func init() {
	register("IAM.Managed Policies", func(r Repos, app *Application) Component {
		return NewIAMPolicies(r.IAM, model.IAMIdentityTypeAll, nil, app)
	})
}
