package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"strconv"
)

type EC2Images struct {
	*ui.Table
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

func (e EC2Images) GetService() string {
	return "EC2"
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
		var amiId, name, status, visibility, platform, architecture, virtualization, rootDevice, blockDevices string
		if v.ImageId != nil {
			amiId = *v.ImageId
		}
		if v.Name != nil {
			name = *v.Name
		}
		status = utils.AutoCase(string(v.State))
		if v.Public != nil {
			visibility = utils.BoolToString(*v.Public, "Public", "Private")
		}
		if v.PlatformDetails != nil {
			platform = *v.PlatformDetails
		}
		architecture = string(v.Architecture)
		virtualization = utils.AutoCase(string(v.VirtualizationType))
		rootDeviceType := utils.AutoCase(string(v.RootDeviceType))
		if v.RootDeviceName != nil {
			rootDevice = fmt.Sprintf("%v (%v)", *v.RootDeviceName, rootDeviceType)
		} else {
			rootDevice = rootDeviceType
		}
		blockDevices = strconv.Itoa(len(v.BlockDeviceMappings))
		data = append(data, []string{
			amiId,
			name,
			status,
			visibility,
			platform,
			architecture,
			virtualization,
			rootDevice,
			blockDevices,
		})
	}
	e.SetData(data)
}
