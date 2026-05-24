package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewGlobalAcceleratorAccelerators(repo *repo.GlobalAccelerator, app *Application) *TableView[model.GlobalAcceleratorAccelerator] {
	return NewTableView(
		"Global Accelerator",
		[]string{"Acclerators"},
		[]string{"NAME", "TYPE", "STATUS", "ENABLED", "ADDRS", "DNS NAME", "DUAL STACK DNS NAME"},
		repo.ListAccelerators,
		func(v model.GlobalAcceleratorAccelerator) []string {
			var enabled string
			if v.Enabled != nil {
				enabled = utils.BoolToString(*v.Enabled, "Yes", "No")
			}
			addrs := 0
			for _, ip := range v.IpSets {
				addrs += len(ip.IpAddresses)
			}
			return []string{
				utils.DerefString(v.Name, ""),
				utils.AutoCase(string(v.IpAddressType)),
				utils.AutoCase(string(v.Status)),
				enabled,
				strconv.Itoa(addrs),
				utils.DerefString(v.DnsName, ""),
				utils.DerefString(v.DualStackDnsName, ""),
			}
		},
		func(v *TableView[model.GlobalAcceleratorAccelerator]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
					Description: "Listeners",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.AcceleratorArn == nil || item.Name == nil {
							return
						}
						app.AddAndSwitch(NewGlobalAcceleratorListeners(repo, *item.Name, *item.AcceleratorArn, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.AcceleratorArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.AcceleratorArn, app))
					},
				},
			}
		},
	)
}
