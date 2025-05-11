package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type EC2SecurityGroups struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2SecurityGroups(repo *repo.EC2, app *Application) *EC2SecurityGroups {
	e := &EC2SecurityGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"VPC ID",
			"INGRESS RULES",
			"EGRESS RULES",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2SecurityGroups) GetLabels() []string {
	return []string{"Security Groups"}
}

func (e EC2SecurityGroups) rulesHandler() {
	sgId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SecurityGroupRules(e.repo, sgId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroups) tagsHandler() {
	sgId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), sgId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
			Description: "Rules",
			Action:      e.rulesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2SecurityGroups) Render() {
	model, err := e.repo.ListSecurityGroups()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.GroupName, ""),
			utils.DerefString(v.GroupId, ""),
			utils.DerefString(v.VpcId, ""),
			strconv.Itoa(len(v.IpPermissions)),
			strconv.Itoa(len(v.IpPermissionsEgress)),
			utils.DerefString(v.Description, ""),
		})
	}
	e.SetData(data)
}
