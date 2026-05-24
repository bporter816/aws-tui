package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewIAMGroups(repo *repo.IAM, userName *string, app *Application) *TableView[model.IAMGroup] {
	var labels []string
	if userName == nil {
		labels = []string{"Groups"}
	} else {
		labels = []string{*userName, "Groups"}
	}
	return NewTableView(
		"IAM",
		labels,
		[]string{"ID", "NAME", "PATH", "CREATED"},
		func() ([]model.IAMGroup, error) {
			return repo.ListGroups(userName)
		},
		func(v model.IAMGroup) []string {
			var created string
			if v.CreateDate != nil {
				created = v.CreateDate.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.GroupId, ""),
				utils.DerefString(v.GroupName, ""),
				utils.DerefString(v.Path, ""),
				created,
			}
		},
		func(v *TableView[model.IAMGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Policies",
					Action: func() {
						groupName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMPolicies(repo, model.IAMIdentityTypeGroup, &groupName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
					Description: "Users",
					Action: func() {
						groupName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewIAMUsers(repo, &groupName, app))
					},
				},
			}
		},
	)
}

func init() {
	register("IAM.Groups", func(r Repos, app *Application) Component {
		return NewIAMGroups(r.IAM, nil, app)
	})
}
