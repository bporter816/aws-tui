package repo

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bporter816/aws-tui/model"
	"time"
)

type Elasticache struct {
	ecClient *ec.Client
}

func NewElasticache(ecClient *ec.Client) *Elasticache {
	return &Elasticache{
		ecClient: ecClient,
	}
}

func (e Elasticache) ListEvents() ([]model.ElasticacheEvent, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -13) // TODO get this closer to the max 14 days
	pg := ec.NewDescribeEventsPaginator(
		e.ecClient,
		&ec.DescribeEventsInput{
			StartTime: aws.Time(oneWeekAgo),
		},
	)
	var events []model.ElasticacheEvent
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheEvent{}, err
		}
		for _, v := range out.Events {
			events = append(events, model.ElasticacheEvent(v))
		}
	}
	return events, nil
}

func (e Elasticache) ListReservedNodes() ([]model.ElasticacheReservedNode, error) {
	pg := ec.NewDescribeReservedCacheNodesPaginator(
		e.ecClient,
		&ec.DescribeReservedCacheNodesInput{},
	)
	var reservations []model.ElasticacheReservedNode
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheReservedNode{}, err
		}
		for _, v := range out.ReservedCacheNodes {
			reservations = append(reservations, model.ElasticacheReservedNode(v))
		}
	}
	return reservations, nil
}

func (e Elasticache) ListSnapshots() ([]model.ElasticacheSnapshot, error) {
	pg := ec.NewDescribeSnapshotsPaginator(
		e.ecClient,
		&ec.DescribeSnapshotsInput{},
	)
	var snapshots []model.ElasticacheSnapshot
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheSnapshot{}, err
		}
		for _, v := range out.Snapshots {
			snapshots = append(snapshots, model.ElasticacheSnapshot(v))
		}
	}
	return snapshots, nil
}

func (e Elasticache) ListParameterGroups() ([]model.ElasticacheParameterGroup, error) {
	pg := ec.NewDescribeCacheParameterGroupsPaginator(
		e.ecClient,
		&ec.DescribeCacheParameterGroupsInput{},
	)
	var parameterGroups []model.ElasticacheParameterGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheParameterGroup{}, err
		}
		for _, v := range out.CacheParameterGroups {
			parameterGroups = append(parameterGroups, model.ElasticacheParameterGroup(v))
		}
	}
	return parameterGroups, nil
}

func (e Elasticache) ListParameters(parameterGroupName string) ([]model.ElasticacheParameter, error) {
	pg := ec.NewDescribeCacheParametersPaginator(
		e.ecClient,
		&ec.DescribeCacheParametersInput{
			CacheParameterGroupName: aws.String(parameterGroupName),
		},
	)
	var parameters []model.ElasticacheParameter
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheParameter{}, err
		}
		for _, v := range out.Parameters {
			parameters = append(parameters, model.ElasticacheParameter(v))
		}
	}
	return parameters, nil
}

func (e Elasticache) ListTags(arn string) (model.Tags, error) {
	out, err := e.ecClient.ListTagsForResource(
		context.TODO(),
		&ec.ListTagsForResourceInput{
			ResourceName: aws.String(arn),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.TagList {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
