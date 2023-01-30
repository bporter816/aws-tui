package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfrontTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type Cloudfront struct {
	client *cloudfront.Client
}

func NewCloudfront(c *cloudfront.Client) *Cloudfront {
	return &Cloudfront{
		client: c,
	}
}

func (c Cloudfront) GetHeaders() []string {
	return []string{
		"ID",
		"DESCRIPTION",
		"STATUS",
		"DOMAIN",
	}
}

func (c *Cloudfront) Render() ([][]string, error) {
	distributionsPaginator := cloudfront.NewListDistributionsPaginator(c.client, &cloudfront.ListDistributionsInput{})
	var distributions []cloudfrontTypes.DistributionSummary
	for distributionsPaginator.HasMorePages() {
		out, err := distributionsPaginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		distributions = append(distributions, out.DistributionList.Items...)
	}

	var ret [][]string
	for _, v := range distributions {
		var comment string
		if v.Comment != nil {
			comment = *v.Comment
		}
		ret = append(ret, []string{
			*v.Id,
			comment,
			*v.Status,
			*v.DomainName,
		})
	}
	return ret, nil
}
