package internal

import (
	"fmt"

	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheUsers(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheUser] {
	return NewTableView(
		"ElastiCache",
		[]string{"Users"},
		[]string{"ID", "NAME", "ACCESS STRING", "STATUS", "AUTH TYPE", "GROUPS"},
		repo.ListUsers,
		func(v model.ElastiCacheUser) []string {
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
			return []string{
				utils.DerefString(v.UserId, ""),
				utils.DerefString(v.UserName, ""),
				utils.DerefString(v.AccessString, ""),
				status,
				authType,
				fmt.Sprintf("%v", len(v.UserGroupIds)),
			}
		},
		func(v *TableView[model.ElastiCacheUser]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ARN == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ARN, app))
					},
				},
			}
		},
	)
}
