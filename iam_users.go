package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type IAMUsers struct {
	*ui.Table
	iamClient *iam.Client
	app       *Application
	groupName string
}

func NewIAMUsers(iamClient *iam.Client, app *Application, groupName string) *IAMUsers {
	i := &IAMUsers{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"PATH",
			"CREATED",
			"PASSWORD LAST USED",
		}, 1, 0),
		iamClient: iamClient,
		app:       app,
		groupName: groupName,
	}
	return i
}

func (i IAMUsers) GetService() string {
	return "IAM"
}

func (i IAMUsers) GetLabels() []string {
	if i.groupName == "" {
		return []string{"Users"}
	} else {
		return []string{i.groupName, "Users"}
	}
}

func (i IAMUsers) accessKeysHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	accessKeysView := NewIAMAccessKeys(i.iamClient, i.app, userName)
	i.app.AddAndSwitch(accessKeysView)
}

func (i IAMUsers) policiesHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	policiesView := NewIAMPolicies(i.iamClient, i.app, IAMIdentityTypeUser, userName)
	i.app.AddAndSwitch(policiesView)
}

func (i IAMUsers) permissionsBoundaryHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	permissionsBoundaryView := NewIAMPolicy(i.iamClient, i.app, IAMIdentityTypeUser, IAMPolicyTypePermissionsBoundary, userName, "", "")
	i.app.AddAndSwitch(permissionsBoundaryView)
}

func (i IAMUsers) groupsHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	groupsView := NewIAMGroups(i.iamClient, i.app, userName)
	i.app.AddAndSwitch(groupsView)
}

func (i IAMUsers) tagsHandler() {
	userName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewIAMUserTags(i.iamClient, i.app, userName)
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
	var users []iamTypes.User

	if i.groupName == "" {
		pg := iam.NewListUsersPaginator(
			i.iamClient,
			&iam.ListUsersInput{},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(context.TODO())
			if err != nil {
				panic(err)
			}
			users = append(users, out.Users...)
		}
	} else {
		pg := iam.NewGetGroupPaginator(
			i.iamClient,
			&iam.GetGroupInput{
				GroupName: aws.String(i.groupName),
			},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(context.TODO())
			if err != nil {
				panic(err)
			}
			users = append(users, out.Users...)
		}
	}

	var data [][]string
	for _, v := range users {
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
			created = v.CreateDate.Format("2006-01-02 15:04:05")
		}
		if v.PasswordLastUsed != nil {
			passwordLastUsed = v.PasswordLastUsed.Format("2006-01-02 15:04:05")
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
