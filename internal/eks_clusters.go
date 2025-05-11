package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type EKSClusters struct {
	*ui.Table
	view.EKS
	repo  *repo.EKS
	app   *Application
	model []model.EKSCluster
}

func NewEKSClusters(repo *repo.EKS, app *Application) *EKSClusters {
	e := &EKSClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"K8S VERSION",
			"OIDC ISSUER",
			"CREATED",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EKSClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (e EKSClusters) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].Arn; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e EKSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *EKSClusters) Render() {
	model, err := e.repo.ListClusters()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var oidcIssuer, created string
		if v.Identity != nil && v.Identity.Oidc != nil && v.Identity.Oidc.Issuer != nil {
			oidcIssuer = *v.Identity.Oidc.Issuer
		}
		if v.CreatedAt != nil {
			created = v.CreatedAt.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			utils.TitleCase(string(v.Status)),
			utils.DerefString(v.Version, ""),
			oidcIssuer,
			created,
		})
	}
	e.SetData(data)
}
