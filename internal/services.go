package internal

import (
	"fmt"
	"sort"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/rivo/tview"
)

//go:generate go run gen.go arg1

type Services struct {
	*ui.Tree
	repos map[string]interface{}
	app   *Application
}

func NewServices(repos map[string]interface{}, app *Application) *Services {
	m := map[string][]string{
		"ACM": {
			"Certificates",
		},
		"ACM PCA": {
			"Certificate Authorities",
		},
		"CloudFront": {
			"Distributions",
			"Functions",
		},
		"CloudWatch": {
			"Log Groups",
		},
		"DynamoDB": {
			"Tables",
		},
		"EBS": {
			"Volumes",
		},
		"EC2": {
			"Instances",
			"Availability Zones",
			"Security Groups",
			"AMIs",
			"Key Pairs",
			"Reserved Instances",
		},
		"ECS": {
			"Clusters",
			"Task Definitions",
		},
		"EKS": {
			"Clusters",
		},
		"ELB": {
			"Load Balancers",
			"Target Groups",
			"Trust Stores",
		},
		"ElastiCache": {
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
		"Global Accelerator": {
			"Accelerators",
		},
		"IAM": {
			"Users",
			"Roles",
			"Groups",
			"Managed Policies",
		},
		"KMS": {
			"Keys",
			"Custom Key Stores",
		},
		"Lambda": {
			"Functions",
		},
		"MQ": {
			"Brokers",
		},
		"MSK": {
			"Clusters",
		},
		"RDS": {
			"Clusters",
			"Global Clusters",
			"Parameter Groups",
			"Subnet Groups",
			"Reserved Instances",
		},
		"Route 53": {
			"Hosted Zones",
			"Health Checks",
		},
		"S3": {
			"Buckets",
		},
		"SNS": {
			"Topics",
		},
		"SQS": {
			"Queues",
		},
		"Secrets Manager": {
			"Secrets",
		},
		"Service Quotas": {
			"Services",
		},
		"Systems Manager": {
			"Parameters",
		},
		"VPC": {
			"VPCs",
			"Subnets",
			"Internet Gateways",
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
	for k := range m {
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
	case "ACM.Certificates":
		item = NewACMCertificates(s.repos["ACM"].(*repo.ACM), s.app)
	case "ACM PCA.Certificate Authorities":
		item = NewACMPCACertificateAuthorities(s.repos["ACM PCA"].(*repo.ACMPCA), s.app)
	case "CloudFront.Distributions":
		item = NewCFDistributions(s.repos["CloudFront"].(*repo.CloudFront), s.app)
	case "CloudFront.Functions":
		item = NewCFFunctions(s.repos["CloudFront"].(*repo.CloudFront), s.app)
	case "CloudWatch.Log Groups":
		item = NewCloudWatchLogGroups(s.repos["CloudWatch"].(*repo.CloudWatch), s.app)
	case "DynamoDB.Tables":
		item = NewDynamoDBTables(s.repos["DynamoDB"].(*repo.DynamoDB), s.app)
	case "EBS.Volumes":
		item = NewEBSVolumes(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Instances":
		item = NewEC2Instances(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Availability Zones":
		item = NewEC2AvailabilityZones(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Security Groups":
		item = NewEC2SecurityGroups(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.AMIs":
		item = NewEC2Images(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Key Pairs":
		item = NewEC2KeyPairs(s.repos["EC2"].(*repo.EC2), s.app)
	case "EC2.Reserved Instances":
		item = NewEC2ReservedInstances(s.repos["EC2"].(*repo.EC2), s.app)
	case "ECS.Clusters":
		item = NewECSClusters(s.repos["ECS"].(*repo.ECS), s.app)
	case "ECS.Task Definitions":
		item = NewECSTaskDefinitions(s.repos["ECS"].(*repo.ECS), s.app)
	case "EKS.Clusters":
		item = NewEKSClusters(s.repos["EKS"].(*repo.EKS), s.app)
	case "ELB.Load Balancers":
		item = NewELBLoadBalancers(s.repos["ELB"].(*repo.ELB), s.app)
	case "ELB.Target Groups":
		item = NewELBTargetGroups(s.repos["ELB"].(*repo.ELB), s.app)
	case "ELB.Trust Stores":
		item = NewELBTrustStores(s.repos["ELB"].(*repo.ELB), s.app)
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
	case "Global Accelerator.Accelerators":
		item = NewGlobalAcceleratorAccelerators(s.repos["Global Accelerator"].(*repo.GlobalAccelerator), s.app)
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
	case "MQ.Brokers":
		item = NewMQBrokers(s.repos["MQ"].(*repo.MQ), s.app)
	case "MSK.Clusters":
		item = NewMSKClusters(s.repos["MSK"].(*repo.MSK), s.app)
	case "RDS.Clusters":
		item = NewRDSClusters(s.repos["RDS"].(*repo.RDS), s.app)
	case "RDS.Global Clusters":
		item = NewRDSGlobalClusters(s.repos["RDS"].(*repo.RDS), s.app)
	case "RDS.Parameter Groups":
		item = NewRDSParameterGroups(s.repos["RDS"].(*repo.RDS), s.app)
	case "RDS.Subnet Groups":
		item = NewRDSSubnetGroups(s.repos["RDS"].(*repo.RDS), s.repos["EC2"].(*repo.EC2), s.app)
	case "RDS.Reserved Instances":
		item = NewRDSReservedInstances(s.repos["RDS"].(*repo.RDS), s.app)
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
	case "Systems Manager.Parameters":
		item = NewSSMParameters(s.repos["SSM"].(*repo.SSM), s.app)
	case "VPC.VPCs":
		item = NewVPCVPCs(s.repos["EC2"].(*repo.EC2), s.app)
	case "VPC.Subnets":
		item = NewVPCSubnets(s.repos["EC2"].(*repo.EC2), []string{}, "", s.app)
	case "VPC.Internet Gateways":
		item = NewVPCInternetGateways(s.repos["EC2"].(*repo.EC2), s.app)
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
