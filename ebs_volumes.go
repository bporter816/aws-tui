package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
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

func (e EBSVolumes) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EBSVolumes) Render() {
	model, err := e.repo.ListVolumes()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, id, volumeType, size, iops, throughput, attachments, encrypted string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		} else {
			name = "-"
		}
		if v.VolumeId != nil {
			id = *v.VolumeId
		}
		volumeType = string(v.VolumeType)
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
			id,
			volumeType,
			size,
			iops,
			throughput,
			attachments,
			encrypted,
		})
	}
	e.SetData(data)
}
