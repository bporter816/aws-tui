package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewEC2AvailabilityZones(repo *repo.EC2, app *Application) *TableView[model.EC2AvailabilityZone] {
	return NewTableView(
		"EC2",
		[]string{"Availability Zones"},
		[]string{"NAME", "ID", "STATE", "MESSAGES"},
		repo.ListAvailabilityZones,
		func(v model.EC2AvailabilityZone) []string {
			var messages string
			if len(v.Messages) > 0 {
				for _, m := range v.Messages {
					if m.Message != nil {
						messages += ", " + *m.Message
					}
				}
				messages = messages[2:]
			} else {
				messages = "-"
			}
			return []string{
				utils.DerefString(v.ZoneName, ""),
				utils.DerefString(v.ZoneId, ""),
				utils.TitleCase(string(v.State)),
				messages,
			}
		},
		nil,
	)
}
