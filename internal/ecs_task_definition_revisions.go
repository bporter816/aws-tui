package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewECSTaskDefinitionRevisions(family string, repo *repo.ECS, app *Application) *TableView[string] {
	return NewTableView(
		"ECS",
		[]string{family, "Revisions"},
		[]string{"REVISION"},
		func() ([]string, error) {
			return repo.ListTaskDefinitionRevisions(family)
		},
		func(v string) []string {
			if a, err := arn.Parse(v); err == nil {
				return []string{utils.GetResourceNameFromArn(a)}
			}
			return []string{v}
		},
		nil,
	)
}
