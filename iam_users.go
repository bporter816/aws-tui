package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
	// "github.com/gdamore/tcell/v2"
)

type IAMUsers struct {
	*ui.Table
	iamClient *iam.Client
	app       *Application
}

func NewIAMUsers(iamClient *iam.Client, app *Application) *IAMUsers {
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
	}
	return i
}

func (i IAMUsers) GetService() string {
	return "IAM"
}

func (i IAMUsers) GetLabels() []string {
	return []string{"Users"}
}

func (i IAMUsers) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMUsers) Render() {
	pg := iam.NewListUsersPaginator(
		i.iamClient,
		&iam.ListUsersInput{},
	)
	var users []iamTypes.User
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		users = append(users, out.Users...)
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
