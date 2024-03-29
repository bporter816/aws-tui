package main

import (
	"errors"

	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type RDSParameters struct {
	*ui.Table
	view.RDS
	repo      *repo.RDS
	app       *Application
	groupName string
	groupType model.RDSParameterGroupType
}

func NewRDSParameters(repo *repo.RDS, app *Application, groupName string, groupType model.RDSParameterGroupType) *RDSParameters {
	r := &RDSParameters{
		Table: ui.NewTable([]string{
			"NAME",
			"VALUE",
			"DATA TYPE",
			"APPLY METHOD",
			"APPLY TYPE",
			"MODIFIABLE",
			"DESCRIPTION",
		}, 1, 1),
		repo:      repo,
		app:       app,
		groupName: groupName,
		groupType: groupType,
	}
	return r
}

func (r RDSParameters) GetLabels() []string {
	return []string{r.groupName, "Parameters"}
}

func (r RDSParameters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r RDSParameters) Render() {
	var parameters []model.RDSParameter
	var err error
	if r.groupType == model.RDSParameterGroupTypeCluster {
		parameters, err = r.repo.ListClusterParameters(r.groupName)
	} else if r.groupType == model.RDSParameterGroupTypeInstance {
		parameters, err = r.repo.ListInstanceParameters(r.groupName)
	} else {
		err = errors.New("param group type must be cluster or instance")
	}
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range parameters {
		var name, value, dataType, applyMethod, applyType, modifiable, desc string
		if v.ParameterName != nil {
			name = *v.ParameterName
		}
		if v.ParameterValue != nil {
			value = *v.ParameterValue
		}
		if v.DataType != nil {
			dataType = *v.DataType
		}
		applyMethod = string(v.ApplyMethod)
		if v.ApplyType != nil {
			applyType = *v.ApplyType
		}
		if v.IsModifiable != nil {
			modifiable = utils.BoolToString(*v.IsModifiable, "Yes", "No")
		}
		if v.Description != nil {
			desc = *v.Description
		}
		data = append(data, []string{
			name,
			value,
			dataType,
			applyMethod,
			applyType,
			modifiable,
			desc,
		})
	}
	r.SetData(data)
}
