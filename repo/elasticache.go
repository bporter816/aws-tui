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

func (e Elasticache) ListClusters() ([]model.ElasticacheCluster, error) {
	// DescribeReplicationGroups doesn't return engine version, so we have to get it from the list of member cluster names
	clusterToEngineVersion := make(map[string]string)

	clustersPg := ec.NewDescribeCacheClustersPaginator(
		e.ecClient,
		&ec.DescribeCacheClustersInput{},
	)
	var clusters []model.ElasticacheCluster
	for clustersPg.HasMorePages() {
		out, err := clustersPg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheCluster{}, err
		}
		for _, v := range out.CacheClusters {
			vCopy := v
			clusters = append(clusters, model.ElasticacheCluster{CacheCluster: &vCopy})
			if v.CacheClusterId != nil && v.EngineVersion != nil && v.ReplicationGroupId != nil {
				clusterToEngineVersion[*v.CacheClusterId] = *v.EngineVersion
			}
		}
	}

	replicationGroupsPg := ec.NewDescribeReplicationGroupsPaginator(
		e.ecClient,
		&ec.DescribeReplicationGroupsInput{},
	)
	for replicationGroupsPg.HasMorePages() {
		out, err := replicationGroupsPg.NextPage(context.TODO())
		if err != nil {
			// TODO more elegantly handle errors
			return clusters, err
		}
		for _, v := range out.ReplicationGroups {
			m := model.ElasticacheCluster{ReplicationGroup: &v}
			if ev, ok := clusterToEngineVersion[v.MemberClusters[0]]; ok {
				m.ReplicationGroupEngineVersion = ev
			}
			clusters = append(clusters, m)
		}
	}
	return clusters, nil
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
