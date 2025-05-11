package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type GlobalAcceleratorAccelerators struct {
	*ui.Table
	view.GlobalAccelerator
	repo  *repo.GlobalAccelerator
	app   *Application
	model []model.GlobalAcceleratorAccelerator
}

func NewGlobalAcceleratorAccelerators(repo *repo.GlobalAccelerator, app *Application) *GlobalAcceleratorAccelerators {
	g := &GlobalAcceleratorAccelerators{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"STATUS",
			"ENABLED",
			"ADDRS",
			"DNS NAME",
			"DUAL STACK DNS NAME",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return g
}

func (g GlobalAcceleratorAccelerators) GetLabels() []string {
	return []string{"Acclerators"}
}

func (g GlobalAcceleratorAccelerators) listenersHandler() {
	row, err := g.GetRowSelection()
	if err != nil {
		return
	}
	if r := g.model[row-1]; r.AcceleratorArn != nil && r.Name != nil {
		listenersView := NewGlobalAcceleratorListeners(g.repo, *r.Name, *r.AcceleratorArn, g.app)
		g.app.AddAndSwitch(listenersView)
	}
}

func (g GlobalAcceleratorAccelerators) tagsHandler() {
	row, err := g.GetRowSelection()
	if err != nil {
		return
	}
	if a := g.model[row-1].AcceleratorArn; a != nil {
		tagsView := NewTags(g.repo, g.GetService(), *a, g.app)
		g.app.AddAndSwitch(tagsView)
	}
}

func (g GlobalAcceleratorAccelerators) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
			Description: "Listeners",
			Action:      g.listenersHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      g.tagsHandler,
		},
	}
}

func (g *GlobalAcceleratorAccelerators) Render() {
	model, err := g.repo.ListAccelerators()
	if err != nil {
		panic(err)
	}
	g.model = model

	var data [][]string
	for _, v := range model {
		var enabled string
		if v.Enabled != nil {
			enabled = utils.BoolToString(*v.Enabled, "Yes", "No")
		}
		addrs := 0
		for _, ip := range v.IpSets {
			addrs += len(ip.IpAddresses)
		}
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			utils.AutoCase(string(v.IpAddressType)),
			utils.AutoCase(string(v.Status)),
			enabled,
			strconv.Itoa(addrs),
			utils.DerefString(v.DnsName, ""),
			utils.DerefString(v.DualStackDnsName, ""),
		})
	}
	g.SetData(data)
}
