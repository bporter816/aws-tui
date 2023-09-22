package main

import (
	"fmt"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type ElasticacheUsers struct {
	*ui.Table
	repo *repo.Elasticache
	app  *Application
}

func NewElasticacheUsers(repo *repo.Elasticache, app *Application) *ElasticacheUsers {
	e := &ElasticacheUsers{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"ACCESS STRING",
			"STATUS",
			"AUTH TYPE",
			"GROUPS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElasticacheUsers) GetService() string {
	return "Elasticache"
}

func (e ElasticacheUsers) GetLabels() []string {
	return []string{"Users"}
}

func (e ElasticacheUsers) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheUsers) Render() {
	model, err := e.repo.ListUsers()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, name, accessString, status, authType, groups string
		if v.UserId != nil {
			id = *v.UserId
		}
		if v.UserName != nil {
			name = *v.UserName
		}
		if v.AccessString != nil {
			accessString = *v.AccessString
		}
		if v.Status != nil {
			status = *v.Status
		}
		if a := v.Authentication; a != nil {
			authType = string(a.Type)
			if a.Type == ecTypes.AuthenticationTypePassword && a.PasswordCount != nil {
				authType += fmt.Sprintf(" (%v)", *a.PasswordCount)
			}
		}
		groups = fmt.Sprintf("%v", len(v.UserGroupIds))
		data = append(data, []string{
			id,
			name,
			accessString,
			status,
			authType,
			groups,
		})
	}
	e.SetData(data)
}
