package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/rivo/tview"
	"sort"
)

type Services struct {
	*ui.Tree
	repos map[string]interface{}
	app   *Application
}

func NewServices(repos map[string]interface{}, app *Application) *Services {
	m := map[string][]string{
		"CloudFront": []string{
			"Distributions",
			"Functions",
		},
		"CloudWatch": []string{
			"Log Groups",
		},
		"DynamoDB": []string{
			"Tables",
		},
		"EC2": []string{
			"Instances",
			"VPCs",
			"Subnets",
			"Availability Zones",
			"Security Groups",
			"Key Pairs",
		},
		"ElastiCache": []string{
			"Clusters",
			"Users",
			"Groups",
			"Parameter Groups",
			"Subnet Groups",
			"Reserved Nodes",
			"Snapshots",
			"Events",
			"Service Updates",
		},
		"ELB": []string{
			"Load Balancers",
			"Target Groups",
		},
		"IAM": []string{
			"Users",
			"Roles",
			"Groups",
			"Managed Policies",
		},
		"KMS": []string{
			"Keys",
			"Custom Key Stores",
		},
		"Lambda": []string{
			"Functions",
		},
		"Route 53": []string{
			"Hosted Zones",
			"Health Checks",
		},
		"S3": []string{
			"Buckets",
		},
		"SNS": []string{
			"Topics",
		},
		"SQS": []string{
			"Queues",
		},
		"Secrets Manager": []string{
			"Secrets",
		},
		"Service Quotas": []string{
			"Services",
		},
	}
	root := tview.NewTreeNode("Services")
	s := &Services{
		Tree:  ui.NewTree(root),
		repos: repos,
		app:   app,
	}
	// sort the keys
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		n := tview.NewTreeNode(k)
		root.AddChild(n)
		for _, view := range v {
			leaf := tview.NewTreeNode(view)
			leaf.SetReference(fmt.Sprintf("%v.%v", k, view))
			n.AddChild(leaf)
		}
		n.CollapseAll()
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s Services) GetService() string {
	return "Services"
}

func (s Services) GetLabels() []string {
	return []string{}
}

func (s Services) selectHandler(n *tview.TreeNode) {
	if n.GetLevel() < 2 {
		if n.IsExpanded() {
			n.Collapse()
		} else {
			n.Expand()
			if len(n.GetChildren()) > 0 {
				s.SetCurrentNode(n.GetChildren()[0])
			}
		}
		return
	}

	view := n.GetReference().(string)
	var item Component
	switch view {
	case "CloudFront.Distributions":
		item = NewCFDistributions(s.repos["CloudFront"].(*repo.CloudFront), s.app)
	case "CloudFront.Functions":
		item = NewCFFunctions(s.repos["CloudFront"].(*repo.CloudFront), s.app)
	case "CloudWatch.Log Groups":
		item = NewCloudWatchLogGroups(s.repos["CloudWatch"].(*repo.CloudWatch), s.app)
	case "DynamoDB.Tables":
		item = NewDynamoDBTables(s.repos["DynamoDB"].(*repo.DynamoDB), s.app)
	case "EC2.Instances":
		item = NewEC2Instances(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.VPCs":
		item = NewEC2VPCs(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Subnets":
		item = NewEC2Subnets(s.repos["EC2"].(*repo.EC2), []string{}, "", s.app)
	case "EC2.Availability Zones":
		item = NewEC2AvailabilityZones(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Security Groups":
		item = NewEC2SecurityGroups(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Key Pairs":
		item = NewEC2KeyPairs(s.repos["EC2"].(*repo.EC2), s.app)
	case "ElastiCache.Clusters":
		item = NewElastiCacheClusters(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Users":
		item = NewElastiCacheUsers(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Groups":
		item = NewElastiCacheGroups(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Parameter Groups":
		item = NewElastiCacheParameterGroups(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Subnet Groups":
		item = NewElastiCacheSubnetGroups(s.repos["ElastiCache"].(*repo.ElastiCache), s.repos["EC2"].(*repo.EC2), s.app)
	case "ElastiCache.Reserved Nodes":
		item = NewElastiCacheReservedCacheNodes(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Snapshots":
		item = NewElastiCacheSnapshots(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Events":
		item = NewElastiCacheEvents(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ElastiCache.Service Updates":
		item = NewElastiCacheServiceUpdates(s.repos["ElastiCache"].(*repo.ElastiCache), s.app)
	case "ELB.Load Balancers":
		item = NewELBLoadBalancers(s.repos["ELB"].(*repo.ELB), s.app)
	case "ELB.Target Groups":
		item = NewELBTargetGroups(s.repos["ELB"].(*repo.ELB), s.app)
	case "IAM.Users":
		item = NewIAMUsers(s.repos["IAM"].(*repo.IAM), nil, s.app)
	case "IAM.Roles":
		item = NewIAMRoles(s.repos["IAM"].(*repo.IAM), s.app)
	case "IAM.Groups":
		item = NewIAMGroups(s.repos["IAM"].(*repo.IAM), nil, s.app)
	case "IAM.Managed Policies":
		item = NewIAMPolicies(s.repos["IAM"].(*repo.IAM), model.IAMIdentityTypeAll, nil, s.app)
	case "KMS.Keys":
		item = NewKmsKeys(s.repos["KMS"].(*repo.KMS), s.app)
	case "KMS.Custom Key Stores":
		item = NewKmsCustomKeyStores(s.repos["KMS"].(*repo.KMS), s.app)
	case "Lambda.Functions":
		item = NewLambdaFunctions(s.repos["Lambda"].(*repo.Lambda), s.app)
	case "Route 53.Hosted Zones":
		item = NewRoute53HostedZones(s.repos["Route 53"].(*repo.Route53), s.app)
	case "Route 53.Health Checks":
		item = NewRoute53HealthChecks(s.repos["Route 53"].(*repo.Route53), s.app)
	case "S3.Buckets":
		item = NewS3Buckets(s.repos["S3"].(*repo.S3), s.app)
	case "SNS.Topics":
		item = NewSNSTopics(s.repos["SNS"].(*repo.SNS), s.app)
	case "SQS.Queues":
		item = NewSQSQueues(s.repos["SQS"].(*repo.SQS), s.app)
	case "Secrets Manager.Secrets":
		item = NewSMSecrets(s.repos["Secrets Manager"].(*repo.SecretsManager), s.app)
	case "Service Quotas.Services":
		item = NewServiceQuotasServices(s.repos["Service Quotas"].(*repo.ServiceQuotas), s.app)
	default:
		panic("unknown service")
	}
	s.app.AddAndSwitch(item)
}

func (s Services) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s Services) Render() {
}
