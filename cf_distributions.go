package main

import (
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strings"
)

type CFDistributions struct {
	*ui.Table
	cfClient *cf.Client
	repo     *repo.Cloudfront
	app      *Application
	model    []model.CloudfrontDistribution
}

func NewCFDistributions(cfClient *cf.Client, repo *repo.Cloudfront, app *Application) *CFDistributions {
	c := &CFDistributions{
		Table: ui.NewTable([]string{
			"ID",
			"DESCRIPTION",
			"TYPE",
			"STATUS",
			"DOMAIN",
			"ALTERNATE DOMAINS",
		}, 1, 0),
		cfClient: cfClient,
		repo:     repo,
		app:      app,
	}
	return c
}

func (c CFDistributions) GetService() string {
	return "Cloudfront"
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
	invalidationsView := NewCFDistributionInvalidations(c.cfClient, id, c.app)
	c.app.AddAndSwitch(invalidationsView)
}

func (c CFDistributions) tagsHandler() {
	row, err := c.GetRowSelection()
	if err != nil {
		return
	}
	if c.model[row-1].ARN == nil {
		return
	}
	tagsView := NewCFTags(c.repo, *c.model[row-1].ARN, c.app)
	c.app.AddAndSwitch(tagsView)
}

func (c CFDistributions) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone),
			Description: "Origins",
			Action:      c.originsHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "Cache Behaviors",
			Action:      c.cacheBehaviorsHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
			Description: "Custom Error Responses",
			Action:      c.customErrorResponsesHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
			Description: "Invalidations",
			Action:      c.invalidationsHandler,
		},
		KeyAction{
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
		var id, comment, distributionType, status, domainName, alternateDomainNames string
		if v.Id != nil {
			id = *v.Id
		}
		if v.Comment != nil {
			comment = *v.Comment
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
			alternateDomainNames = strings.Join(v.Aliases.Items, ", ")
		}
		data = append(data, []string{
			id,
			comment,
			distributionType,
			status,
			domainName,
			alternateDomainNames,
		})
	}
	c.SetData(data)
}
