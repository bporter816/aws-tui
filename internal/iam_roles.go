package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewIAMRoles(repo *repo.IAM, app *Application) *TableView[model.IAMRole] {
	return NewTableView(
		"IAM",
		[]string{"Roles"},
		[]string{"ID", "NAME", "PATH", "MAX SESSION", "CREATED", "DESCRIPTION"},
		repo.ListRoles,
		func(v model.IAMRole) []string {
			var maxSession, created string
			if v.MaxSessionDuration != nil {
				maxSession = strconv.Itoa(int(*v.MaxSessionDuration))
			}
			if v.CreateDate != nil {
				created = v.CreateDate.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.RoleId, ""),
				utils.DerefString(v.RoleName, ""),
				utils.DerefString(v.Path, ""),
				maxSession,
				created,
				utils.DerefString(v.Description, ""),
			}
		},
		func(v *TableView[model.IAMRole]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Policies",
					Action: func() {
						roleName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMPolicies(repo, model.IAMIdentityTypeRole, &roleName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
					Description: "Assume Role Policy",
					Action: func() {
						roleName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMPolicy(repo, model.IAMIdentityTypeRole, model.IAMPolicyTypeAssumeRolePolicy, roleName, "", "", app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
					Description: "Permissions Boundary",
					Action: func() {
						roleName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMPolicy(repo, model.IAMIdentityTypeRole, model.IAMPolicyTypePermissionsBoundary, roleName, "", "", app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						roleName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), "role:"+roleName, app))
					},
				},
			}
		},
	)
}

func init() {
	register("IAM.Roles", func(r Repos, app *Application) Component {
		return NewIAMRoles(r.IAM, app)
	})
}
