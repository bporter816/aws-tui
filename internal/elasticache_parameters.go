package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewElastiCacheParameters(repo *repo.ElastiCache, parameterGroupName string, app *Application) *TableView[model.ElastiCacheParameter] {
	return NewTableView(
		"ElastiCache",
		[]string{parameterGroupName, "Parameters"},
		[]string{"NAME", "ALLOWED VALUES", "VALUE", "TYPE", "MODIFIABLE", "CHANGE TYPE", "SOURCE", "DESCRIPTION"},
		func(ctx context.Context) ([]model.ElastiCacheParameter, error) {
			return repo.ListParameters(ctx, parameterGroupName)
		},
		func(v model.ElastiCacheParameter) []string {
			var isModifiable string
			if v.IsModifiable != nil {
				isModifiable = utils.BoolToString(*v.IsModifiable, "Yes", "No")
			}
			return []string{
				utils.DerefString(v.ParameterName, ""),
				utils.DerefString(v.AllowedValues, ""),
				utils.DerefString(v.ParameterValue, ""),
				utils.DerefString(v.DataType, ""),
				isModifiable,
				string(v.ChangeType),
				utils.DerefString(v.Source, ""),
				utils.DerefString(v.Description, ""),
			}
		},
		nil,
	)
}
