package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/bporter816/aws-tui/internal/model"
)

type ECS struct {
	ecsClient *ecs.Client
}

func NewECS(ecsClient *ecs.Client) *ECS {
	return &ECS{
		ecsClient: ecsClient,
	}
}

// Internal function to get all cluster arns
func (e ECS) listClusterArns(ctx context.Context) ([]string, error) {
	pg := ecs.NewListClustersPaginator(
		e.ecsClient,
		&ecs.ListClustersInput{},
	)
	var arns []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		arns = append(arns, out.ClusterArns...)
	}
	return arns, nil
}

func (e ECS) ListClusters(ctx context.Context) ([]model.ECSCluster, error) {
	arns, err := e.listClusterArns(ctx)
	if err != nil {
		return []model.ECSCluster{}, err
	}
	var clusters []model.ECSCluster
	out, err := e.ecsClient.DescribeClusters(
		ctx,
		&ecs.DescribeClustersInput{
			Clusters: arns,
		},
	)
	if err != nil {
		return []model.ECSCluster{}, err
	}
	for _, v := range out.Clusters {
		// TODO handle failures
		clusters = append(clusters, model.ECSCluster(v))
	}
	return clusters, nil
}

// Internal function to get all service arns
func (e ECS) listServiceArns(ctx context.Context, clusterName string) ([]string, error) {
	pg := ecs.NewListServicesPaginator(
		e.ecsClient,
		&ecs.ListServicesInput{
			Cluster: aws.String(clusterName),
		},
	)
	var arns []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		arns = append(arns, out.ServiceArns...)
	}
	return arns, nil
}

func (e ECS) ListServices(ctx context.Context, clusterName string) ([]model.ECSService, error) {
	arns, err := e.listServiceArns(ctx, clusterName)
	if err != nil {
		return []model.ECSService{}, err
	}
	if len(arns) == 0 {
		return []model.ECSService{}, nil
	}
	var services []model.ECSService
	out, err := e.ecsClient.DescribeServices(
		ctx,
		&ecs.DescribeServicesInput{
			Cluster:  aws.String(clusterName),
			Services: arns,
		},
	)
	if err != nil {
		return []model.ECSService{}, err
	}
	for _, v := range out.Services {
		// TODO handle failures
		services = append(services, model.ECSService(v))
	}
	return services, nil
}

// Internal function to get task arns
func (e ECS) listTaskArns(ctx context.Context, clusterName string, serviceName string) ([]string, error) {
	input := &ecs.ListTasksInput{
		Cluster: aws.String(clusterName),
	}
	if len(serviceName) > 0 {
		input.ServiceName = aws.String(serviceName)
	}
	pg := ecs.NewListTasksPaginator(e.ecsClient, input)
	var arns []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		arns = append(arns, out.TaskArns...)
	}
	return arns, nil
}

func (e ECS) ListTasks(ctx context.Context, clusterName string, serviceName string) ([]model.ECSTask, error) {
	arns, err := e.listTaskArns(ctx, clusterName, serviceName)
	if err != nil {
		return []model.ECSTask{}, err
	}
	if len(arns) == 0 {
		return []model.ECSTask{}, nil
	}
	var tasks []model.ECSTask
	out, err := e.ecsClient.DescribeTasks(
		ctx,
		&ecs.DescribeTasksInput{
			Cluster: aws.String(clusterName),
			Tasks:   arns,
		},
	)
	if err != nil {
		return []model.ECSTask{}, err
	}
	for _, v := range out.Tasks {
		// TODO handle failures
		tasks = append(tasks, model.ECSTask(v))
	}
	return tasks, nil
}

func (e ECS) ListTaskDefinitions(ctx context.Context) ([]string, error) {
	pg := ecs.NewListTaskDefinitionFamiliesPaginator(
		e.ecsClient,
		&ecs.ListTaskDefinitionFamiliesInput{},
	)
	var families []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		families = append(families, out.Families...)
	}
	return families, nil
}

func (e ECS) ListTaskDefinitionRevisions(ctx context.Context, family string) ([]string, error) {
	pg := ecs.NewListTaskDefinitionsPaginator(
		e.ecsClient,
		&ecs.ListTaskDefinitionsInput{
			FamilyPrefix: aws.String(family),
		},
	)
	var revisions []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		revisions = append(revisions, out.TaskDefinitionArns...)
	}
	return revisions, nil
}

func (e ECS) ListTags(ctx context.Context, resourceId string) (model.Tags, error) {
	out, err := e.ecsClient.ListTagsForResource(
		ctx,
		&ecs.ListTagsForResourceInput{
			ResourceArn: aws.String(resourceId),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
