package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cwLogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/bporter816/aws-tui/model"
)

type CloudWatch struct {
	cwLogsClient *cwLogs.Client
}

func NewCloudWatch(cwLogsClient *cwLogs.Client) *CloudWatch {
	return &CloudWatch{
		cwLogsClient: cwLogsClient,
	}
}

func (c CloudWatch) ListLogGroups() ([]model.CloudWatchLogGroup, error) {
	pg := cwLogs.NewDescribeLogGroupsPaginator(
		c.cwLogsClient,
		&cwLogs.DescribeLogGroupsInput{},
	)
	var logGroups []model.CloudWatchLogGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.CloudWatchLogGroup{}, err
		}
		for _, v := range out.LogGroups {
			logGroups = append(logGroups, model.CloudWatchLogGroup(v))
		}
	}
	return logGroups, nil
}

func (c CloudWatch) ListTags(resourceArn string) (model.Tags, error) {
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
