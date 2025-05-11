package internal

import (
	"fmt"

	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheUsers struct {
	*ui.Table
	view.ElastiCache
	repo  *repo.ElastiCache
	app   *Application
	model []model.ElastiCacheUser
}

func NewElastiCacheUsers(repo *repo.ElastiCache, app *Application) *ElastiCacheUsers {
	e := &ElastiCacheUsers{
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

func (e ElastiCacheUsers) GetLabels() []string {
	return []string{"Users"}
}

func (e ElastiCacheUsers) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].ARN; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheUsers) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElastiCacheUsers) Render() {
	model, err := e.repo.ListUsers()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var status, authType string
		if v.Status != nil {
			status = utils.TitleCase(*v.Status)
		}
		if a := v.Authentication; a != nil {
			authType = utils.AutoCase(string(a.Type))
			if a.Type == ecTypes.AuthenticationTypePassword && a.PasswordCount != nil {
				authType += fmt.Sprintf(" (%v)", *a.PasswordCount)
			}
		}
		data = append(data, []string{
			utils.DerefString(v.UserId, ""),
			utils.DerefString(v.UserName, ""),
			utils.DerefString(v.AccessString, ""),
			status,
			authType,
			fmt.Sprintf("%v", len(v.UserGroupIds)),
		})
	}
	e.SetData(data)
}
