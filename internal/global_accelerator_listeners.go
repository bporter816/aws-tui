package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewGlobalAcceleratorListeners(repo *repo.GlobalAccelerator, acceleratorName string, acceleratorArn string, app *Application) *TableView[model.GlobalAcceleratorListener] {
	return NewTableView(
		"Global Accelerator",
		[]string{acceleratorName, "Listeners"},
		[]string{"PROTOCOL", "PORT RANGES", "CLIENT AFFINITY"},
		func() ([]model.GlobalAcceleratorListener, error) {
			return repo.ListListeners(acceleratorArn)
		},
		func(v model.GlobalAcceleratorListener) []string {
			return []string{
				utils.AutoCase(string(v.Protocol)),
				utils.FormatGlobalAcceleratorPortRanges(v.PortRanges),
				utils.AutoCase(string(v.ClientAffinity)),
			}
		},
		nil,
	)
}
