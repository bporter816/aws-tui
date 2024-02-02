package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type SSMParameters struct {
	*ui.Table
	view.SSM
	repo *repo.SSM
	app  *Application
}

func NewSSMParameters(repo *repo.SSM, app *Application) *SSMParameters {
	s := &SSMParameters{
		Table: ui.NewTable([]string{
			"NAME",
			"TIER",
			"TYPE",
			"DATA TYPE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return s
}

func (s SSMParameters) GetLabels() []string {
	return []string{"Parameters"}
}

func (s SSMParameters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SSMParameters) Render() {
	model, err := s.repo.ListParameters()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, tier, parameterType, dataType, version, policies string
		if v.Name != nil {
			name = *v.Name
		}
		tier = string(v.Tier)
		parameterType = string(v.Type)
		if v.DataType != nil {
			dataType = *v.DataType
		}
		version = strconv.FormatInt(v.Version, 10)
		policies = strconv.Itoa(len(v.Policies))
		data = append(data, []string{
			name,
			tier,
			parameterType,
			dataType,
			version,
			policies,
		})
	}
	s.SetData(data)
}
