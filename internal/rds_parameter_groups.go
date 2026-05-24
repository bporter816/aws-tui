package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type rdsParameterGroup struct {
	name        string
	groupType   string
	family      string
	description string
	arn         string
}

func NewRDSParameterGroups(repo *repo.RDS, app *Application) *TableView[rdsParameterGroup] {
	return NewTableView(
		"RDS",
		[]string{"Parameter Groups"},
		[]string{"NAME", "TYPE", "FAMILY", "DESCRIPTION"},
		func() ([]rdsParameterGroup, error) {
			clusterGroups, err := repo.ListClusterParameterGroups()
			if err != nil {
				return nil, err
			}
			instanceGroups, err := repo.ListInstanceParameterGroups()
			if err != nil {
				return nil, err
			}
			var groups []rdsParameterGroup
			for _, v := range clusterGroups {
				groups = append(groups, rdsParameterGroup{
					name:        utils.DerefString(v.DBClusterParameterGroupName, ""),
					groupType:   "Cluster",
					family:      utils.DerefString(v.DBParameterGroupFamily, ""),
					description: utils.DerefString(v.Description, ""),
					arn:         v.Arn(),
				})
			}
			for _, v := range instanceGroups {
				groups = append(groups, rdsParameterGroup{
					name:        utils.DerefString(v.DBParameterGroupName, ""),
					groupType:   "Instance",
					family:      utils.DerefString(v.DBParameterGroupFamily, ""),
					description: utils.DerefString(v.Description, ""),
					arn:         v.Arn(),
				})
			}
			return groups, nil
		},
		func(v rdsParameterGroup) []string {
			return []string{
				v.name,
				v.groupType,
				v.family,
				v.description,
			}
		},
		func(v *TableView[rdsParameterGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Parameters",
					Action: func() {
						name, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						groupType, err := v.GetColSelection("TYPE")
						if err != nil {
							return
						}
						var parameterGroupType model.RDSParameterGroupType
						if groupType == "Cluster" {
							parameterGroupType = model.RDSParameterGroupTypeCluster
						} else {
							parameterGroupType = model.RDSParameterGroupTypeInstance
						}
						app.AddAndSwitch(NewRDSParameters(repo, app, name, parameterGroupType))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), item.arn, app))
					},
				},
			}
		},
	)
}
