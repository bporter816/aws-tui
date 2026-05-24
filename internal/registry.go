package internal

import (
	"sort"
	"strings"

	"github.com/bporter816/aws-tui/internal/repo"
)

// Repos holds all AWS service repository instances.
type Repos struct {
	ACM               *repo.ACM
	ACMPCA            *repo.ACMPCA
	CloudFront        *repo.CloudFront
	CloudWatch        *repo.CloudWatch
	DynamoDB          *repo.DynamoDB
	EC2               *repo.EC2
	ECS               *repo.ECS
	EKS               *repo.EKS
	ELB               *repo.ELB
	ElastiCache       *repo.ElastiCache
	GlobalAccelerator *repo.GlobalAccelerator
	IAM               *repo.IAM
	KMS               *repo.KMS
	Lambda            *repo.Lambda
	MQ                *repo.MQ
	MSK               *repo.MSK
	RDS               *repo.RDS
	Route53           *repo.Route53
	S3                *repo.S3
	SecretsManager    *repo.SecretsManager
	ServiceQuotas     *repo.ServiceQuotas
	SNS               *repo.SNS
	SQS               *repo.SQS
	SSM               *repo.SSM
}

type componentFactory func(Repos, *Application) Component

var registry = map[string]componentFactory{}

func register(key string, f componentFactory) {
	registry[key] = f
}

// serviceTree derives the services navigation tree from registered components.
// Returns a map of service name -> alphabetically sorted list of view names.
func serviceTree() map[string][]string {
	tree := map[string][]string{}
	for key := range registry {
		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			continue
		}
		tree[parts[0]] = append(tree[parts[0]], parts[1])
	}
	for service := range tree {
		sort.Strings(tree[service])
	}
	return tree
}
