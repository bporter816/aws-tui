package internal

import (
	"errors"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewRDSParameters(repo *repo.RDS, app *Application, groupName string, groupType model.RDSParameterGroupType) *TableView[model.RDSParameter] {
	return NewTableView(
		"RDS",
		[]string{groupName, "Parameters"},
		[]string{"NAME", "VALUE", "DATA TYPE", "APPLY METHOD", "APPLY TYPE", "MODIFIABLE", "DESCRIPTION"},
		func() ([]model.RDSParameter, error) {
			if groupType == model.RDSParameterGroupTypeCluster {
				return repo.ListClusterParameters(groupName)
			} else if groupType == model.RDSParameterGroupTypeInstance {
				return repo.ListInstanceParameters(groupName)
			}
			return nil, errors.New("param group type must be cluster or instance")
		},
		func(v model.RDSParameter) []string {
			var modifiable string
			if v.IsModifiable != nil {
				modifiable = utils.BoolToString(*v.IsModifiable, "Yes", "No")
			}
			return []string{
				utils.DerefString(v.ParameterName, ""),
				utils.DerefString(v.ParameterValue, ""),
				utils.DerefString(v.DataType, ""),
				string(v.ApplyMethod),
				utils.DerefString(v.ApplyType, ""),
				modifiable,
				utils.DerefString(v.Description, ""),
			}
		},
		nil,
	)
}
