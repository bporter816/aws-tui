package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cwLogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/bporter816/aws-tui/model"
)

type Cloudwatch struct {
	cwLogsClient *cwLogs.Client
}

func NewCloudwatch(cwLogsClient *cwLogs.Client) *Cloudwatch {
	return &Cloudwatch{
		cwLogsClient: cwLogsClient,
	}
}

func (c Cloudwatch) ListLogGroups() ([]model.CloudwatchLogGroup, error) {
	pg := cwLogs.NewDescribeLogGroupsPaginator(
		c.cwLogsClient,
		&cwLogs.DescribeLogGroupsInput{},
	)
	var logGroups []model.CloudwatchLogGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.CloudwatchLogGroup{}, err
		}
		for _, v := range out.LogGroups {
			logGroups = append(logGroups, model.CloudwatchLogGroup(v))
		}
	}
	return logGroups, nil
}

func (c Cloudwatch) ListTags(resourceArn string) (model.Tags, error) {
	out, err := c.cwLogsClient.ListTagsForResource(
		context.TODO(),
		&cwLogs.ListTagsForResourceInput{
			ResourceArn: aws.String(resourceArn),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for k, v := range out.Tags {
		tags = append(tags, model.Tag{Key: k, Value: v})
	}
	return tags, nil
}
