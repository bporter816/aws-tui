package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type IAMUsers struct {
	*ui.Table
	view.IAM
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
	tagsView := NewTags(i.repo, i.GetService(), "user:"+userName, i.app)
	i.app.AddAndSwitch(tagsView)
}

func (i IAMUsers) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Access Keys",
			Action:      i.accessKeysHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Policies",
			Action:      i.policiesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
			Description: "Permissions Boundary",
			Action:      i.permissionsBoundaryHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
			Description: "Groups",
			Action:      i.groupsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var created, passwordLastUsed string
		if v.CreateDate != nil {
			created = v.CreateDate.Format(utils.DefaultTimeFormat)
		}
		if v.PasswordLastUsed != nil {
			passwordLastUsed = v.PasswordLastUsed.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.UserId, ""),
			utils.DerefString(v.UserName, ""),
			utils.DerefString(v.Path, ""),
			created,
			passwordLastUsed,
		})
	}
	i.SetData(data)
}
