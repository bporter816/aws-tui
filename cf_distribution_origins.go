package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type CFDistributionOrigins struct {
	*ui.Table
	view.CloudFront
	repo           *repo.CloudFront
	distributionId string
	app            *Application
}

func NewCFDistributionOrigins(repo *repo.CloudFront, distributionId string, app *Application) *CFDistributionOrigins {
	c := &CFDistributionOrigins{
		Table: ui.NewTable([]string{
			"NAME",
			"DOMAIN NAME",
			"PATH",
			"TYPE",
		}, 1, 0),
		repo:           repo,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionOrigins) GetLabels() []string {
	return []string{c.distributionId, "Origins"}
}

func (c CFDistributionOrigins) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionOrigins) Render() {
	model, err := c.repo.GetDistributionOrigins(c.distributionId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, domainName, originPath, originType string
		if v.Id != nil {
			id = *v.Id
		}
		if v.DomainName != nil {
			domainName = *v.DomainName
		}
		if v.OriginPath != nil {
			originPath = *v.OriginPath
		}
		if v.S3OriginConfig != nil {
			originType = "S3"
		} else if v.CustomOriginConfig != nil {
			originType = "Custom origin"
		}
		data = append(data, []string{
			id,
			domainName,
			originPath,
			originType,
		})
	}
	c.SetData(data)
}
