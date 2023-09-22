package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type IAMUsers struct {
	*ui.Table
	repo      *repo.IAM
	groupName *string
	app       *Application
}

func NewIAMUsers(repo *repo.IAM, groupName *string, app *Application) *IAMUsers {
	i := &IAMUsers{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"PATH",
			"CREATED",
			"PASSWORD LAST USED",
		}, 1, 0),
		repo:      repo,
		groupName: groupName,
		app:       app,
	}
	return i
}

func (i IAMUsers) GetService() string {
	return "IAM"
}

func (i IAMUsers) GetLabels() []string {
	if i.groupName == nil {
		return []string{"Users"}
	} else {
		return []string{*i.groupName, "Users"}
	}
}

func (i IAMUsers) accessKeysHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	accessKeysView := NewIAMAccessKeys(i.repo, userName, i.app)
	i.app.AddAndSwitch(accessKeysView)
}

func (i IAMUsers) policiesHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	policiesView := NewIAMPolicies(i.repo, model.IAMIdentityTypeUser, &userName, i.app)
	i.app.AddAndSwitch(policiesView)
}

func (i IAMUsers) permissionsBoundaryHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	permissionsBoundaryView := NewIAMPolicy(i.repo, model.IAMIdentityTypeUser, model.IAMPolicyTypePermissionsBoundary, userName, "", "", i.app)
	i.app.AddAndSwitch(permissionsBoundaryView)
}

func (i IAMUsers) groupsHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	groupsView := NewIAMGroups(i.repo, &userName, i.app)
	i.app.AddAndSwitch(groupsView)
}

func (i IAMUsers) tagsHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewIAMUserTags(i.repo, userName, i.app)
	i.app.AddAndSwitch(tagsView)
}

func (i IAMUsers) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Access Keys",
			Action:      i.accessKeysHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Policies",
			Action:      i.policiesHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
			Description: "Permissions Boundary",
			Action:      i.permissionsBoundaryHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
			Description: "Groups",
			Action:      i.groupsHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      i.tagsHandler,
		},
	}
}

func (i IAMUsers) Render() {
	model, err := i.repo.ListUsers(i.groupName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var userId, userName, path, created, passwordLastUsed string
		if v.UserId != nil {
			userId = *v.UserId
		}
		if v.UserName != nil {
			userName = *v.UserName
		}
		if v.Path != nil {
			path = *v.Path
		}
		if v.CreateDate != nil {
			created = v.CreateDate.Format(utils.DefaultTimeFormat)
		}
		if v.PasswordLastUsed != nil {
			passwordLastUsed = v.PasswordLastUsed.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			userId,
			userName,
			path,
			created,
			passwordLastUsed,
		})
	}
	i.SetData(data)
}
