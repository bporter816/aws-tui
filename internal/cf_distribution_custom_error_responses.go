package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
)

func NewCFDistributionCustomErrorResponses(repo *repo.CloudFront, distributionId string, app *Application) *TableView[model.CloudFrontDistributionCustomErrorResponse] {
	return NewTableView(
		"CloudFront",
		[]string{distributionId, "Custom Error Responses"},
		[]string{"ERROR CODE", "RESPONSE CODE", "RESPONSE PAGE PATH", "MIN TTL SECONDS"},
		func() ([]model.CloudFrontDistributionCustomErrorResponse, error) {
			return repo.GetDistributionCustomErrorResponses(distributionId)
		},
		func(v model.CloudFrontDistributionCustomErrorResponse) []string {
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
			return []string{
				errorCode,
				responseCode,
				responsePagePath,
				minTTL,
			}
		},
		nil,
	)
}
