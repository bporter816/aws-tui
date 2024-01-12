package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type EC2Volumes struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2Volumes(repo *repo.EC2, app *Application) *EC2Volumes {
	e := &EC2Volumes{
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

func (e EC2Volumes) GetLabels() []string {
	return []string{"Volumes"}
}

func (e EC2Volumes) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2Volumes) Render() {
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
