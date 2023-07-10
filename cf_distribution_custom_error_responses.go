package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"strconv"
)

type CFDistributionCustomErrorResponses struct {
	*Table
	cfClient       *cf.Client
	distributionId string
	app            *Application
}

func NewCFDistributionCustomErrorResponses(cfClient *cf.Client, distributionId string, app *Application) *CFDistributionCustomErrorResponses {
	c := &CFDistributionCustomErrorResponses{
		Table: NewTable([]string{
			"ERROR CODE",
			"RESPONSE CODE",
			"RESPONSE PAGE PATH",
			"MIN TTL SECONDS",
		}, 1, 0),
		cfClient:       cfClient,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionCustomErrorResponses) GetService() string {
	return "Cloudfront"
}

func (c CFDistributionCustomErrorResponses) GetLabels() []string {
	return []string{c.distributionId, "Custom Error Responses"}
}

func (c CFDistributionCustomErrorResponses) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionCustomErrorResponses) Render() {
	out, err := c.cfClient.GetDistributionConfig(
		context.TODO(),
		&cf.GetDistributionConfigInput{
			Id: aws.String(c.distributionId),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.DistributionConfig != nil && out.DistributionConfig.CustomErrorResponses != nil {
		for _, v := range out.DistributionConfig.CustomErrorResponses.Items {
			var errorCode, responseCode, responsePagePath, minTTL string
			if v.ErrorCode != nil {
				errorCode = strconv.Itoa(int(*v.ErrorCode))
			}
			if v.ResponseCode != nil {
				responseCode = *v.ResponseCode
			}
			if v.ResponsePagePath != nil {
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
	}
	c.SetData(data)
}