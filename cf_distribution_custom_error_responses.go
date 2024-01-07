package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type CFDistributionCustomErrorResponses struct {
	*ui.Table
	view.CloudFront
	repo           *repo.CloudFront
	distributionId string
	app            *Application
}

func NewCFDistributionCustomErrorResponses(repo *repo.CloudFront, distributionId string, app *Application) *CFDistributionCustomErrorResponses {
	c := &CFDistributionCustomErrorResponses{
		Table: ui.NewTable([]string{
			"ERROR CODE",
			"RESPONSE CODE",
			"RESPONSE PAGE PATH",
			"MIN TTL SECONDS",
		}, 1, 0),
		repo:           repo,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionCustomErrorResponses) GetLabels() []string {
	return []string{c.distributionId, "Custom Error Responses"}
}

func (c CFDistributionCustomErrorResponses) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionCustomErrorResponses) Render() {
	model, err := c.repo.GetDistributionCustomErrorResponses(c.distributionId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		errorCode, responseCode, responsePagePath, minTTL := "-", "-", "-", "-"
		if v.ErrorCode != nil {
			errorCode = strconv.Itoa(int(*v.ErrorCode))
		}
		if v.ResponseCode != nil && *v.ResponseCode != "" {
			responseCode = *v.ResponseCode
		}
		if v.ResponsePagePath != nil && *v.ResponsePagePath != "" {
			responsePagePath = *v.ResponsePagePath
		}
		if v.ErrorCachingMinTTL != nil {
			minTTL = strconv.FormatInt(*v.ErrorCachingMinTTL, 10)
		}
		data = append(data, []string{
			errorCode,
			responseCode,
			responsePagePath,
			minTTL,
		})
	}
	c.SetData(data)
}
