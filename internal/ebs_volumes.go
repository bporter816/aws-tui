package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type EBSVolumes struct {
	*ui.Table
	view.EBS
	repo *repo.EC2
	app  *Application
}

func NewEBSVolumes(repo *repo.EC2, app *Application) *EBSVolumes {
	e := &EBSVolumes{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"TYPE",
			"SIZE",
			"IOPS",
			"THROUGHPUT",
			"ATTACHMENTS",
			"ENCRYPTED",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EBSVolumes) GetLabels() []string {
	return []string{"Volumes"}
}

func (e EBSVolumes) tagsHandler() {
	id, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), id, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EBSVolumes) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EBSVolumes) Render() {
	model, err := e.repo.ListVolumes(nil)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
		data = append(data, []string{
			name,
			utils.DerefString(v.VolumeId, ""),
			string(v.VolumeType),
			size,
			iops,
			throughput,
			attachments,
			encrypted,
		})
	}
	e.SetData(data)
}
