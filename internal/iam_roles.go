package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type IAMRoles struct {
	*ui.Table
	view.IAM
	repo *repo.IAM
	app  *Application
}

func NewIAMRoles(repo *repo.IAM, app *Application) *IAMRoles {
	i := &IAMRoles{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"PATH",
			"MAX SESSION",
			"CREATED",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return i
}

func (i IAMRoles) GetLabels() []string {
	return []string{"Roles"}
}

func (i IAMRoles) policiesHandler() {
	roleName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	policiesView := NewIAMPolicies(i.repo, model.IAMIdentityTypeRole, &roleName, i.app)
	i.app.AddAndSwitch(policiesView)
}

func (i IAMRoles) assumeRolePolicyHandler() {
	roleName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	assumeRolePolicyView := NewIAMPolicy(i.repo, model.IAMIdentityTypeRole, model.IAMPolicyTypeAssumeRolePolicy, roleName, "", "", i.app)
	i.app.AddAndSwitch(assumeRolePolicyView)
}

func (i IAMRoles) permissionsBoundaryHandler() {
	roleName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	permissionsBoundaryView := NewIAMPolicy(i.repo, model.IAMIdentityTypeRole, model.IAMPolicyTypePermissionsBoundary, roleName, "", "", i.app)
	i.app.AddAndSwitch(permissionsBoundaryView)
}

func (i IAMRoles) tagsHandler() {
	roleName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewTags(i.repo, i.GetService(), "role:"+roleName, i.app)
	i.app.AddAndSwitch(tagsView)
}

func (i IAMRoles) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Policies",
			Action:      i.policiesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Assume Role Policy",
			Action:      i.assumeRolePolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
			Description: "Permissions Boundary",
			Action:      i.permissionsBoundaryHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      i.tagsHandler,
		},
	}
}

func (i IAMRoles) Render() {
	model, err := i.repo.ListRoles()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var maxSession, created string
		if v.MaxSessionDuration != nil {
			maxSession = strconv.Itoa(int(*v.MaxSessionDuration))
		}
		if v.CreateDate != nil {
			created = v.CreateDate.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.RoleId, ""),
			utils.DerefString(v.RoleName, ""),
			utils.DerefString(v.Path, ""),
			maxSession,
			created,
			utils.DerefString(v.Description, ""),
		})
	}
	i.SetData(data)
}
