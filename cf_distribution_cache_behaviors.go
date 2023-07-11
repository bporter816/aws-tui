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
			"CACHE POLICY",
			"ORIGIN REQUEST POLICY",
			"RESPONSE HEADERS POLICY",
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
				cachePolicyId, originRequestPolicyId, responseHeadersPolicyId := "-", "-", "-"
				if v.PathPattern != nil {
					pathPattern = *v.PathPattern
				}
				if v.TargetOriginId != nil {
					origin = *v.TargetOriginId
				}
				viewerProtocolPolicy = viewerProtocolPolicyToString(v.ViewerProtocolPolicy)
				if v.CachePolicyId != nil {
					cachePolicyId = *v.CachePolicyId
				}
				if v.OriginRequestPolicyId != nil {
					originRequestPolicyId = *v.OriginRequestPolicyId
				}
				if v.ResponseHeadersPolicyId != nil {
					responseHeadersPolicyId = *v.ResponseHeadersPolicyId
				}
				data = append(data, []string{
					pathPattern,
					origin,
					viewerProtocolPolicy,
					cachePolicyId,
					originRequestPolicyId,
					responseHeadersPolicyId,
				})
			}
		}
		if d := out.DistributionConfig.DefaultCacheBehavior; d != nil {
			var origin, viewerProtocolPolicy string
			cachePolicyId, originRequestPolicyId, responseHeadersPolicyId := "-", "-", "-"
			if d.TargetOriginId != nil {
				origin = *d.TargetOriginId
			}
			viewerProtocolPolicy = viewerProtocolPolicyToString(d.ViewerProtocolPolicy)
			if d.CachePolicyId != nil {
				cachePolicyId = *d.CachePolicyId
			}
			if d.OriginRequestPolicyId != nil {
				originRequestPolicyId = *d.OriginRequestPolicyId
			}
			if d.ResponseHeadersPolicyId != nil {
				responseHeadersPolicyId = *d.ResponseHeadersPolicyId
			}
			data = append(data, []string{
				"Default (*)",
				origin,
				viewerProtocolPolicy,
				cachePolicyId,
				originRequestPolicyId,
				responseHeadersPolicyId,
			})
		}
	}
	c.SetData(data)
}
