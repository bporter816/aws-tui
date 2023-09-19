package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type IAMGroups struct {
	*ui.Table
	repo      *repo.IAM
	iamClient *iam.Client
	userName  string
	app       *Application
}

func NewIAMGroups(repo *repo.IAM, iamClient *iam.Client, userName string, app *Application) *IAMGroups {
	i := &IAMGroups{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"PATH",
			"CREATED",
		}, 1, 0),
		repo:      repo,
		iamClient: iamClient,
		userName:  userName,
		app:       app,
	}
	return i
}

func (i IAMGroups) GetService() string {
	return "IAM"
}

func (i IAMGroups) GetLabels() []string {
	if i.userName == "" {
		return []string{"Groups"}
	} else {
		return []string{i.userName, "Groups"}
	}
}

func (i IAMGroups) policiesHandler() {
	groupName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	policiesView := NewIAMPolicies(i.repo, i.iamClient, model.IAMIdentityTypeGroup, groupName, i.app)
	i.app.AddAndSwitch(policiesView)
}

func (i IAMGroups) usersHandler() {
	groupName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	usersView := NewIAMUsers(i.repo, i.iamClient, groupName, i.app)
	i.app.AddAndSwitch(usersView)
}

func (i IAMGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Policies",
			Action:      i.policiesHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
			Description: "Users",
			Action:      i.usersHandler,
		},
	}
}

func (i IAMGroups) Render() {
	var groups []iamTypes.Group

	// IAM has different APIs for the user-scoped and global groups lists
	if i.userName == "" {
		pg := iam.NewListGroupsPaginator(
			i.iamClient,
			&iam.ListGroupsInput{},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(context.TODO())
			if err != nil {
				panic(err)
			}
			groups = append(groups, out.Groups...)
		}
	} else {
		pg := iam.NewListGroupsForUserPaginator(
			i.iamClient,
			&iam.ListGroupsForUserInput{
				UserName: aws.String(i.userName),
			},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(context.TODO())
			if err != nil {
				panic(err)
			}
			groups = append(groups, out.Groups...)
		}
	}

	var data [][]string
	for _, v := range groups {
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
