package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

type CFDistributionCacheBehaviors struct {
	*Table
	cfClient       *cf.Client
	distributionId string
	app            *Application
}

func NewCFDistributionCacheBehaviors(cfClient *cf.Client, distributionId string, app *Application) *CFDistributionCacheBehaviors {
	c := &CFDistributionCacheBehaviors{
		Table: NewTable([]string{
			"PATH",
			"ORIGIN",
			"VIEWER PROTOCOL POLICY",
		}, 1, 0),
		cfClient:       cfClient,
		distributionId: distributionId,
		app:            app,
	}
	return c
}

func (c CFDistributionCacheBehaviors) GetService() string {
	return "Cloudfront"
}

func (c CFDistributionCacheBehaviors) GetLabels() []string {
	return []string{c.distributionId, "Cache Behaviors"}
}

func (c CFDistributionCacheBehaviors) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFDistributionCacheBehaviors) Render() {
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
	if out.DistributionConfig != nil {
		if out.DistributionConfig.CacheBehaviors != nil {
			for _, v := range out.DistributionConfig.CacheBehaviors.Items {
				var pathPattern, origin, viewerProtocolPolicy string
				if v.PathPattern != nil {
					pathPattern = *v.PathPattern
				}
				if v.TargetOriginId != nil {
					origin = *v.TargetOriginId
				}
				viewerProtocolPolicy = viewerProtocolPolicyToString(v.ViewerProtocolPolicy)
				data = append(data, []string{
					pathPattern,
					origin,
					viewerProtocolPolicy,
				})
			}
		}
		if d := out.DistributionConfig.DefaultCacheBehavior; d != nil {
			var defaultOrigin, defaultViewerProtocolPolicy string
			if d.TargetOriginId != nil {
				defaultOrigin = *d.TargetOriginId
			}
			defaultViewerProtocolPolicy = viewerProtocolPolicyToString(d.ViewerProtocolPolicy)
			data = append(data, []string{
				"Default (*)",
				defaultOrigin,
				defaultViewerProtocolPolicy,
			})
		}
	}
	c.SetData(data)
}
