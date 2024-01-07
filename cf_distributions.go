package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type CFDistributions struct {
	*ui.Table
	view.CloudFront
	repo  *repo.CloudFront
	app   *Application
	model []model.CloudFrontDistribution
}

func NewCFDistributions(repo *repo.CloudFront, app *Application) *CFDistributions {
	c := &CFDistributions{
		Table: ui.NewTable([]string{
			"ID",
			"TYPE",
			"STATUS",
			"DOMAIN",
			"ALTERNATE DOMAINS",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return c
}

func (c CFDistributions) GetLabels() []string {
	return []string{"Distributions"}
}

func (c CFDistributions) originsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	originsView := NewCFDistributionOrigins(c.repo, id, c.app)
	c.app.AddAndSwitch(originsView)
}

func (c CFDistributions) cacheBehaviorsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	cacheBehaviorsView := NewCFDistributionCacheBehaviors(c.repo, id, c.app)
	c.app.AddAndSwitch(cacheBehaviorsView)
}

func (c CFDistributions) customErrorResponsesHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	customErrorResponsesView := NewCFDistributionCustomErrorResponses(c.repo, id, c.app)
	c.app.AddAndSwitch(customErrorResponsesView)
}

func (c CFDistributions) invalidationsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	invalidationsView := NewCFDistributionInvalidations(c.repo, id, c.app)
	c.app.AddAndSwitch(invalidationsView)
}

func (c CFDistributions) tagsHandler() {
	row, err := c.GetRowSelection()
	if err != nil {
		return
	}
	if arn := c.model[row-1].ARN; arn != nil {
		tagsView := NewCFTags(c.repo, *arn, c.app)
		c.app.AddAndSwitch(tagsView)
	}
}

func (c CFDistributions) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone),
			Description: "Origins",
			Action:      c.originsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "Cache Behaviors",
			Action:      c.cacheBehaviorsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
			Description: "Custom Error Responses",
			Action:      c.customErrorResponsesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
			Description: "Invalidations",
			Action:      c.invalidationsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      c.tagsHandler,
		},
	}
}

func (c *CFDistributions) Render() {
	model, err := c.repo.ListDistributions()
	if err != nil {
		panic(err)
	}
	c.model = model

	var data [][]string
	for _, v := range model {
		var id, distributionType, status, domainName, alternateDomainNames, description string
		if v.Id != nil {
			id = *v.Id
		}
		if v.Staging != nil {
			if *v.Staging {
				distributionType = "Staging"
			} else {
				distributionType = "Production"
			}
		}
		if v.Status != nil {
			status = *v.Status
		}
		if v.DomainName != nil {
			domainName = *v.DomainName
		}
		if v.Aliases != nil && len(v.Aliases.Items) > 0 {
			alternateDomainNames = utils.TruncateStrings(v.Aliases.Items, 1)
		}
		if v.Comment != nil {
			description = *v.Comment
		}
		data = append(data, []string{
			id,
			distributionType,
			status,
			domainName,
			alternateDomainNames,
			description,
		})
	}
	c.SetData(data)
}
