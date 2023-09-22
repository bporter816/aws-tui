package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type IAMRoles struct {
	*ui.Table
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

func (i IAMRoles) GetService() string {
	return "IAM"
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
	tagsView := NewIAMRoleTags(i.repo, roleName, i.app)
	i.app.AddAndSwitch(tagsView)
}

func (i IAMRoles) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Policies",
			Action:      i.policiesHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Assume Role Policy",
			Action:      i.assumeRolePolicyHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
			Description: "Permissions Boundary",
			Action:      i.permissionsBoundaryHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
		var roleId, roleName, path, maxSession, created, description string
		if v.RoleId != nil {
			roleId = *v.RoleId
		}
		if v.RoleName != nil {
			roleName = *v.RoleName
		}
		if v.Path != nil {
			path = *v.Path
		}
		if v.MaxSessionDuration != nil {
			maxSession = strconv.Itoa(int(*v.MaxSessionDuration))
		}
		if v.CreateDate != nil {
			created = v.CreateDate.Format(utils.DefaultTimeFormat)
		}
		if v.Description != nil {
			description = *v.Description
		}
		data = append(data, []string{
			roleId,
			roleName,
			path,
			maxSession,
			created,
			description,
		})
	}
	i.SetData(data)
}
