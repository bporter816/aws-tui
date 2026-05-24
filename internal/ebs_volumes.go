package internal

import (
	"context"
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEBSVolumes(repo *repo.EC2, app *Application) *TableView[model.EC2Volume] {
	return NewTableView(
		"EBS",
		[]string{"Volumes"},
		[]string{"NAME", "ID", "TYPE", "SIZE", "IOPS", "THROUGHPUT", "ATTACHMENTS", "ENCRYPTED"},
		func(ctx context.Context) ([]model.EC2Volume, error) {
			return repo.ListVolumes(ctx, nil)
		},
		func(v model.EC2Volume) []string {
			var name, size, iops, throughput, attachments, encrypted string
			if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
				name = n
			}
			if v.Size != nil {
				size = strconv.Itoa(int(*v.Size)) + " GiB"
			}
			if v.Iops != nil {
				iops = strconv.Itoa(int(*v.Iops))
			}
			if v.Throughput != nil {
				throughput = strconv.Itoa(int(*v.Throughput))
			}
			attachments = strconv.Itoa(len(v.Attachments))
			if v.Encrypted != nil {
				encrypted = utils.BoolToString(*v.Encrypted, "Yes", "No")
			}
			return []string{
				name,
				utils.DerefString(v.VolumeId, ""),
				string(v.VolumeType),
				size,
				iops,
				throughput,
				attachments,
				encrypted,
			}
		},
		func(v *TableView[model.EC2Volume]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), id, app))
					},
				},
			}
		},
	)
}

func init() {
	register("EBS.Volumes", func(r Repos, app *Application) Component {
		return NewEBSVolumes(r.EC2, app)
	})
}
