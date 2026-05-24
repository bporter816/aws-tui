package internal

import (
	"fmt"
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewEC2Images(repo *repo.EC2, app *Application) *TableView[model.EC2Image] {
	return NewTableView(
		"EC2",
		[]string{"Images"},
		[]string{"AMI ID", "NAME", "STATUS", "VISIBILITY", "PLATFORM", "ARCHITECTURE", "VIRTUALIZATION", "ROOT DEVICE", "BLOCK DEVICES"},
		repo.ListImages,
		func(v model.EC2Image) []string {
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
			return []string{
				utils.DerefString(v.ImageId, ""),
				utils.DerefString(v.Name, ""),
				utils.AutoCase(string(v.State)),
				visibility,
				utils.DerefString(v.PlatformDetails, ""),
				string(v.Architecture),
				utils.AutoCase(string(v.VirtualizationType)),
				rootDevice,
				strconv.Itoa(len(v.BlockDeviceMappings)),
			}
		},
		nil,
	)
}

func init() {
	register("EC2.AMIs", func(r Repos, app *Application) Component {
		return NewEC2Images(r.EC2, app)
	})
}
