package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewCFDistributions(repo *repo.CloudFront, app *Application) *TableView[model.CloudFrontDistribution] {
	return NewTableView(
		"CloudFront",
		[]string{"Distributions"},
		[]string{"ID", "TYPE", "STATUS", "DOMAIN", "ALTERNATE DOMAINS", "DESCRIPTION"},
		repo.ListDistributions,
		func(v model.CloudFrontDistribution) []string {
			var distributionType, alternateDomainNames string
			if v.Staging != nil {
				if *v.Staging {
					distributionType = "Staging"
				} else {
					distributionType = "Production"
				}
			}
			if v.Aliases != nil && len(v.Aliases.Items) > 0 {
				alternateDomainNames = utils.TruncateStrings(v.Aliases.Items, 1)
			}
			return []string{
				utils.DerefString(v.Id, ""),
				distributionType,
				utils.DerefString(v.Status, ""),
				utils.DerefString(v.DomainName, ""),
				alternateDomainNames,
				utils.DerefString(v.Comment, ""),
			}
		},
		func(v *TableView[model.CloudFrontDistribution]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone),
					Description: "Origins",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewCFDistributionOrigins(repo, id, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
					Description: "Cache Behaviors",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewCFDistributionCacheBehaviors(repo, id, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
					Description: "Custom Error Responses",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewCFDistributionCustomErrorResponses(repo, id, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
					Description: "Invalidations",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewCFDistributionInvalidations(repo, id, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ARN == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ARN, app))
					},
				},
			}
		},
	)
}
