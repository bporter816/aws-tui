package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewIAMUsers(repo *repo.IAM, groupName *string, app *Application) *TableView[model.IAMUser] {
	var labels []string
	if groupName == nil {
		labels = []string{"Users"}
	} else {
		labels = []string{*groupName, "Users"}
	}
	return NewTableView(
		"IAM",
		labels,
		[]string{"ID", "NAME", "PATH", "CREATED", "PASSWORD LAST USED"},
		func(ctx context.Context) ([]model.IAMUser, error) {
			return repo.ListUsers(ctx, groupName)
		},
		func(v model.IAMUser) []string {
			var created, passwordLastUsed string
			if v.CreateDate != nil {
				created = v.CreateDate.Format(utils.DefaultTimeFormat)
			}
			if v.PasswordLastUsed != nil {
				passwordLastUsed = v.PasswordLastUsed.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.UserId, ""),
				utils.DerefString(v.UserName, ""),
				utils.DerefString(v.Path, ""),
				created,
				passwordLastUsed,
			}
		},
		func(v *TableView[model.IAMUser]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
					Description: "Access Keys",
					Action: func() {
						userName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMAccessKeys(repo, userName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Policies",
					Action: func() {
						userName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMPolicies(repo, model.IAMIdentityTypeUser, &userName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
					Description: "Permissions Boundary",
					Action: func() {
						userName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMPolicy(repo, model.IAMIdentityTypeUser, model.IAMPolicyTypePermissionsBoundary, userName, "", "", app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
					Description: "Groups",
					Action: func() {
						userName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMGroups(repo, &userName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						userName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), "user:"+userName, app))
					},
				},
			}
		},
	)
}

func init() {
	register("IAM.Users", func(r Repos, app *Application) Component {
		return NewIAMUsers(r.IAM, nil, app)
	})
}
