package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type GlobalAcceleratorListeners struct {
	*ui.Table
	view.GlobalAccelerator
	repo            *repo.GlobalAccelerator
	acceleratorName string
	acceleratorArn  string
	app             *Application
	model           []model.GlobalAcceleratorListener
}

func NewGlobalAcceleratorListeners(repo *repo.GlobalAccelerator, acceleratorName string, acceleratorArn string, app *Application) *GlobalAcceleratorListeners {
	g := &GlobalAcceleratorListeners{
		Table: ui.NewTable([]string{
			"PROTOCOL",
			"PORT RANGES",
			"CLIENT AFFINITY",
		}, 1, 0),
		repo:            repo,
		acceleratorName: acceleratorName,
		acceleratorArn:  acceleratorArn,
		app:             app,
	}
	return g
}

func (g GlobalAcceleratorListeners) GetLabels() []string {
	return []string{g.acceleratorName, "Listeners"}
}

func (g GlobalAcceleratorListeners) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (g *GlobalAcceleratorListeners) Render() {
	model, err := g.repo.ListListeners(g.acceleratorArn)
	if err != nil {
		panic(err)
	}
	g.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.AutoCase(string(v.Protocol)),
			utils.FormatGlobalAcceleratorPortRanges(v.PortRanges),
			utils.AutoCase(string(v.ClientAffinity)),
		})
	}
	g.SetData(data)
}
