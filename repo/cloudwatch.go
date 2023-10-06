package repo

import (
	"context"
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
