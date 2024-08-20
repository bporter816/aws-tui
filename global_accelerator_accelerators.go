package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"strconv"
)

type GlobalAcceleratorAccelerators struct {
	*ui.Table
	view.GlobalAccelerator
	repo *repo.GlobalAccelerator
	app  *Application
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

func (g GlobalAcceleratorAccelerators) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (g GlobalAcceleratorAccelerators) Render() {
	model, err := g.repo.ListAccelerators()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, acceleratorType, status, enabled, dnsName, dualStackDnsName string
		if v.Name != nil {
			name = *v.Name
		}
		acceleratorType = utils.AutoCase(string(v.IpAddressType))
		status = utils.AutoCase(string(v.Status))
		if v.Enabled != nil {
			enabled = utils.BoolToString(*v.Enabled, "Yes", "No")
		}
		addrs := 0
		for _, ip := range v.IpSets {
			addrs += len(ip.IpAddresses)
		}
		if v.DnsName != nil {
			dnsName = *v.DnsName
		}
		if v.DualStackDnsName != nil {
			dualStackDnsName = *v.DualStackDnsName
		}
		data = append(data, []string{
			name,
			acceleratorType,
			status,
			enabled,
			strconv.Itoa(addrs),
			dnsName,
			dualStackDnsName,
		})
	}
	g.SetData(data)
}
