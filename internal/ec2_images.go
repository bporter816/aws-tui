package internal

import (
	"fmt"
	"strconv"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type EC2Images struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2Images(repo *repo.EC2, app *Application) *EC2Images {
	e := &EC2Images{
		Table: ui.NewTable([]string{
			"AMI ID",
			"NAME",
			"STATUS",
			"VISIBILITY",
			"PLATFORM",
			"ARCHITECTURE",
			"VIRTUALIZATION",
			"ROOT DEVICE",
			"BLOCK DEVICES",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2Images) GetLabels() []string {
	return []string{"Images"}
}

func (e EC2Images) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2Images) Render() {
	model, err := e.repo.ListImages()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var visibility, rootDevice string
		if v.Public != nil {
			visibility = utils.BoolToString(*v.Public, "Public", "Private")
		}
		rootDeviceType := utils.AutoCase(string(v.RootDeviceType))
		if v.RootDeviceName != nil {
			rootDevice = fmt.Sprintf("%v (%v)", *v.RootDeviceName, rootDeviceType)
		} else {
			rootDevice = rootDeviceType
		}
		data = append(data, []string{
			utils.DerefString(v.ImageId, ""),
			utils.DerefString(v.Name, ""),
			utils.AutoCase(string(v.State)),
			visibility,
			utils.DerefString(v.PlatformDetails, ""),
			string(v.Architecture),
			utils.AutoCase(string(v.VirtualizationType)),
			rootDevice,
			strconv.Itoa(len(v.BlockDeviceMappings)),
		})
	}
	e.SetData(data)
}
