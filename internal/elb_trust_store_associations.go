package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type ELBTrustStoreAssociations struct {
	*ui.Table
	view.ELB
	repo *repo.ELB
	a    arn.ARN
	app  *Application
}

func NewELBTrustStoreAssociations(repo *repo.ELB, a arn.ARN, app *Application) *ELBTrustStoreAssociations {
	e := &ELBTrustStoreAssociations{
		Table: ui.NewTable([]string{
			"ARN",
		}, 1, 0),
		repo: repo,
		a:    a,
		app:  app,
	}
	return e
}

func (e ELBTrustStoreAssociations) GetLabels() []string {
	return []string{e.a.Resource, "Associations"}
}

func (e ELBTrustStoreAssociations) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTrustStoreAssociations) Render() {
	model, err := e.repo.ListTrustStoreAssociations(e.a)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.ResourceArn, ""),
		})
	}
	e.SetData(data)
}
