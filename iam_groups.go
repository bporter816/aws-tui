package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type IAMGroups struct {
	*ui.Table
	repo     *repo.IAM
	userName *string
	app      *Application
}

func NewIAMGroups(repo *repo.IAM, userName *string, app *Application) *IAMGroups {
	i := &IAMGroups{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"PATH",
			"CREATED",
		}, 1, 0),
		repo:     repo,
		userName: userName,
		app:      app,
	}
	return i
}

func (i IAMGroups) GetService() string {
	return "IAM"
}

func (i IAMGroups) GetLabels() []string {
	if i.userName == nil {
		return []string{"Groups"}
	} else {
		return []string{*i.userName, "Groups"}
	}
}

func (i IAMGroups) policiesHandler() {
	groupName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	policiesView := NewIAMPolicies(i.repo, model.IAMIdentityTypeGroup, &groupName, i.app)
	i.app.AddAndSwitch(policiesView)
}

func (i IAMGroups) usersHandler() {
	groupName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	usersView := NewIAMUsers(i.repo, &groupName, i.app)
	i.app.AddAndSwitch(usersView)
}

func (i IAMGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Policies",
			Action:      i.policiesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
			Description: "Users",
			Action:      i.usersHandler,
		},
	}
}

func (i IAMGroups) Render() {
	model, err := i.repo.ListGroups(i.userName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var groupId, groupName, path, created string
		if v.GroupId != nil {
			groupId = *v.GroupId
		}
		if v.GroupName != nil {
			groupName = *v.GroupName
		}
		if v.Path != nil {
			path = *v.Path
		}
		if v.CreateDate != nil {
			created = v.CreateDate.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			groupId,
			groupName,
			path,
			created,
		})
	}
	i.SetData(data)
}
