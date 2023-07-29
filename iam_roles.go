package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type IAMRoles struct {
	*ui.Table
	iamClient *iam.Client
	app       *Application
}

func NewIAMRoles(iamClient *iam.Client, app *Application) *IAMRoles {
	i := &IAMRoles{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"PATH",
			"MAX SESSION",
			"CREATED",
			"DESCRIPTION",
		}, 1, 0),
		iamClient: iamClient,
		app:       app,
	}
	return i
}

func (i IAMRoles) GetService() string {
	return "IAM"
}

func (i IAMRoles) GetLabels() []string {
	return []string{"Roles"}
}

func (i IAMRoles) tagsHandler() {
	roleName, err := i.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewIAMRoleTags(i.iamClient, i.app, roleName)
	i.app.AddAndSwitch(tagsView)
}

func (i IAMRoles) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      i.tagsHandler,
		},
	}
}

func (i IAMRoles) Render() {
	var roles []iamTypes.Role
	pg := iam.NewListRolesPaginator(
		i.iamClient,
		&iam.ListRolesInput{},
	)
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		roles = append(roles, out.Roles...)
	}

	var data [][]string
	for _, v := range roles {
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
			created = v.CreateDate.Format("2006-01-02 15:04:05")
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
