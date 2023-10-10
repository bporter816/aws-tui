package repo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	cw "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/bporter816/aws-tui/model"
	"time"
)

type ElastiCache struct {
	ecClient *ec.Client
	cwClient *cw.Client
}

func NewElastiCache(ecClient *ec.Client, cwClient *cw.Client) *ElastiCache {
	return &ElastiCache{
		ecClient: ecClient,
		cwClient: cwClient,
	}
}

func (e ElastiCache) ListClusters() ([]model.ElastiCacheCluster, error) {
	// DescribeReplicationGroups doesn't return engine version, so we have to get it from the list of member cluster names
	clusterToEngineVersion := make(map[string]string)

	clustersPg := ec.NewDescribeCacheClustersPaginator(
		e.ecClient,
		&ec.DescribeCacheClustersInput{},
	)
	var clusters []model.ElastiCacheCluster
	for clustersPg.HasMorePages() {
		out, err := clustersPg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheCluster{}, err
		}
		for _, v := range out.CacheClusters {
			// skip clusters in replication groups as those are retrieved from DescribeReplicationGroups
			if v.ReplicationGroupId != nil {
				continue
			}
			vCopy := v
			clusters = append(clusters, model.ElastiCacheCluster{CacheCluster: &vCopy})
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
			vCopy := v
			m := model.ElastiCacheCluster{ReplicationGroup: &vCopy}
			if ev, ok := clusterToEngineVersion[v.MemberClusters[0]]; ok {
				m.ReplicationGroupEngineVersion = ev
			}
			clusters = append(clusters, m)
		}
	}
	return clusters, nil
}

func (e ElastiCache) ListEvents() ([]model.ElastiCacheEvent, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -13) // TODO get this closer to the max 14 days
	pg := ec.NewDescribeEventsPaginator(
		e.ecClient,
		&ec.DescribeEventsInput{
			StartTime: aws.Time(oneWeekAgo),
		},
	)
	var events []model.ElastiCacheEvent
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheEvent{}, err
		}
		for _, v := range out.Events {
			events = append(events, model.ElastiCacheEvent(v))
		}
	}
	return events, nil
}

func (e ElastiCache) ListReservedNodes() ([]model.ElastiCacheReservedNode, error) {
	pg := ec.NewDescribeReservedCacheNodesPaginator(
		e.ecClient,
		&ec.DescribeReservedCacheNodesInput{},
	)
	var reservations []model.ElastiCacheReservedNode
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheReservedNode{}, err
		}
		for _, v := range out.ReservedCacheNodes {
			reservations = append(reservations, model.ElastiCacheReservedNode(v))
		}
	}
	return reservations, nil
}

func (e ElastiCache) ListSnapshots() ([]model.ElastiCacheSnapshot, error) {
	pg := ec.NewDescribeSnapshotsPaginator(
		e.ecClient,
		&ec.DescribeSnapshotsInput{},
	)
	var snapshots []model.ElastiCacheSnapshot
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheSnapshot{}, err
		}
		for _, v := range out.Snapshots {
			snapshots = append(snapshots, model.ElastiCacheSnapshot(v))
		}
	}
	return snapshots, nil
}

func (e ElastiCache) ListSubnetGroups() ([]model.ElastiCacheSubnetGroup, error) {
	pg := ec.NewDescribeCacheSubnetGroupsPaginator(
		e.ecClient,
		&ec.DescribeCacheSubnetGroupsInput{},
	)
	var subnetGroups []model.ElastiCacheSubnetGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheSubnetGroup{}, err
		}
		for _, v := range out.CacheSubnetGroups {
			subnetGroups = append(subnetGroups, model.ElastiCacheSubnetGroup(v))
		}
	}
	return subnetGroups, nil
}

func (e ElastiCache) ListParameterGroups() ([]model.ElastiCacheParameterGroup, error) {
	pg := ec.NewDescribeCacheParameterGroupsPaginator(
		e.ecClient,
		&ec.DescribeCacheParameterGroupsInput{},
	)
	var parameterGroups []model.ElastiCacheParameterGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheParameterGroup{}, err
		}
		for _, v := range out.CacheParameterGroups {
			parameterGroups = append(parameterGroups, model.ElastiCacheParameterGroup(v))
		}
	}
	return parameterGroups, nil
}

func (e ElastiCache) ListParameters(parameterGroupName string) ([]model.ElastiCacheParameter, error) {
	pg := ec.NewDescribeCacheParametersPaginator(
		e.ecClient,
		&ec.DescribeCacheParametersInput{
			CacheParameterGroupName: aws.String(parameterGroupName),
		},
	)
	var parameters []model.ElastiCacheParameter
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheParameter{}, err
		}
		for _, v := range out.Parameters {
			parameters = append(parameters, model.ElastiCacheParameter(v))
		}
	}
	return parameters, nil
}

func (e ElastiCache) ListUsers() ([]model.ElastiCacheUser, error) {
	pg := ec.NewDescribeUsersPaginator(
		e.ecClient,
		&ec.DescribeUsersInput{},
	)
	var users []model.ElastiCacheUser
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheUser{}, err
		}
		for _, v := range out.Users {
			users = append(users, model.ElastiCacheUser(v))
		}
	}
	return users, nil
}

func (e ElastiCache) ListGroups() ([]model.ElastiCacheGroup, error) {
	pg := ec.NewDescribeUserGroupsPaginator(
		e.ecClient,
		&ec.DescribeUserGroupsInput{},
	)
	var users []model.ElastiCacheGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheGroup{}, err
		}
		for _, v := range out.UserGroups {
			users = append(users, model.ElastiCacheGroup(v))
		}
	}
	return users, nil
}

func (e ElastiCache) ListServiceUpdates() ([]model.ElastiCacheServiceUpdate, error) {
	pg := ec.NewDescribeServiceUpdatesPaginator(
		e.ecClient,
		&ec.DescribeServiceUpdatesInput{},
	)
	var serviceUpdates []model.ElastiCacheServiceUpdate
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheServiceUpdate{}, err
		}
		for _, v := range out.ServiceUpdates {
			serviceUpdates = append(serviceUpdates, model.ElastiCacheServiceUpdate(v))
		}
	}
	return serviceUpdates, nil
}

func (e ElastiCache) ListUpdateActions(
	cacheClusterIds []string,
	replicationGroupIds []string,
	serviceUpdateName string,
) ([]model.ElastiCacheUpdateAction, error) {
	hasCacheClusters, hasReplicationGroups, hasServiceUpdate := 0, 0, 0
	if len(cacheClusterIds) > 0 {
		hasCacheClusters = 1
	}
	if len(replicationGroupIds) > 0 {
		hasReplicationGroups = 1
	}
	if len(serviceUpdateName) > 0 {
		hasServiceUpdate = 1
	}
	if hasCacheClusters+hasReplicationGroups+hasServiceUpdate != 1 {
		return []model.ElastiCacheUpdateAction{}, errors.New("must specify either cacheClusterIds, replicationGroupIds, or serviceUpdateName")
	}

	var input ec.DescribeUpdateActionsInput
	if len(cacheClusterIds) > 0 {
		input.CacheClusterIds = cacheClusterIds
	} else if len(replicationGroupIds) > 0 {
		input.ReplicationGroupIds = replicationGroupIds
	} else {
		input.ServiceUpdateName = aws.String(serviceUpdateName)
	}
	pg := ec.NewDescribeUpdateActionsPaginator(
		e.ecClient,
		&input,
	)
	var updateActions []model.ElastiCacheUpdateAction
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElastiCacheUpdateAction{}, err
		}
		for _, v := range out.UpdateActions {
			updateActions = append(updateActions, model.ElastiCacheUpdateAction(v))
		}
	}
	return updateActions, nil
}

func (e ElastiCache) ListTags(arn string) (model.Tags, error) {
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
