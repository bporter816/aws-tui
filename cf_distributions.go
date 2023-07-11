package main

import (
	"context"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/gdamore/tcell/v2"
	"strings"
)

type CFDistributions struct {
	*Table
	cfClient *cf.Client
	app      *Application
	arns     []string
}

func NewCFDistributions(cfClient *cf.Client, app *Application) *CFDistributions {
	c := &CFDistributions{
		Table: NewTable([]string{
			"ID",
			"DESCRIPTION",
			"TYPE",
			"STATUS",
			"DOMAIN",
			"ALTERNATE DOMAINS",
		}, 1, 0),
		cfClient: cfClient,
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
	originsView := NewCFDistributionOrigins(c.cfClient, id, c.app)
	c.app.AddAndSwitch(originsView)
}

func (c CFDistributions) cacheBehaviorsHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	cacheBehaviorsView := NewCFDistributionCacheBehaviors(c.cfClient, id, c.app)
	c.app.AddAndSwitch(cacheBehaviorsView)
}

func (c CFDistributions) customErrorResponsesHandler() {
	id, err := c.GetColSelection("ID")
	if err != nil {
		return
	}
	customErrorResponsesView := NewCFDistributionCustomErrorResponses(c.cfClient, id, c.app)
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
	tagsView := NewCFTags(c.cfClient, c.arns[row-1], c.app)
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
	pg := cf.NewListDistributionsPaginator(
		c.cfClient,
		&cf.ListDistributionsInput{},
	)
	var distributions []cfTypes.DistributionSummary
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		distributions = append(distributions, out.DistributionList.Items...)
	}

	var data [][]string
	c.arns = make([]string, len(distributions))
	for i, v := range distributions {
		c.arns[i] = *v.ARN
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
