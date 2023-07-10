package main

import (
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func viewerProtocolPolicyToString(p cfTypes.ViewerProtocolPolicy) string {
	switch p {
	case cfTypes.ViewerProtocolPolicyAllowAll:
		return "Allow all"
	case cfTypes.ViewerProtocolPolicyHttpsOnly:
		return "HTTPS only"
	case cfTypes.ViewerProtocolPolicyRedirectToHttps:
		return "Redirect to HTTPS"
	default:
		return "<unknown>"
	}
}
