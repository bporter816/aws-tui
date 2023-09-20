package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/rivo/tview"
	"sort"
)

type Services struct {
	*ui.Tree
	clients map[string]interface{}
	repos   map[string]interface{}
	app     *Application
}

func NewServices(clients map[string]interface{}, repos map[string]interface{}, app *Application) *Services {
	m := map[string][]string{
		"Cloudfront": []string{
			"Distributions",
			"Functions",
		},
		"DynamoDB": []string{
			"Tables",
		},
		"EC2": []string{
			"Instances",
			"VPCs",
			"Security Groups",
			"Key Pairs",
		},
		"Elasticache": []string{
			"Clusters",
			"Events",
			"Parameter Groups",
			"Reserved Nodes",
			"Snapshots",
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
		"Route 53": []string{
			"Hosted Zones",
			"Health Checks",
		},
		"S3": []string{
			"Buckets",
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
		Tree:    ui.NewTree(root),
		clients: clients,
		repos:   repos,
		app:     app,
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

	s.Root.Walk(func(node, parent *tview.TreeNode) bool {
		// Skip non-leaf nodes but continue traversing.
		// We have to skip the root because it has the same name as Service Quotas "Services".
		if node.GetLevel() < 2 {
			return true
		}

		if node.GetText() == n.GetText() {
			view := fmt.Sprintf("%v.%v", parent.GetText(), node.GetText())
			var item Component
			switch view {
			case "Cloudfront.Distributions":
				item = NewCFDistributions(s.repos["Cloudfront"].(*repo.Cloudfront), s.app)
			case "Cloudfront.Functions":
				item = NewCFFunctions(s.repos["Cloudfront"].(*repo.Cloudfront), s.app)
			case "DynamoDB.Tables":
				item = NewDynamoDBTables(s.repos["DynamoDB"].(*repo.DynamoDB), s.app)
			case "EC2.Instances":
				item = NewEC2Instances(s.repos["EC2"].(*repo.EC2), s.app)
			case "EC2.VPCs":
				item = NewEC2VPCs(s.repos["EC2"].(*repo.EC2), s.app)
			case "EC2.Security Groups":
				item = NewEC2SecurityGroups(s.repos["EC2"].(*repo.EC2), s.app)
			case "EC2.Key Pairs":
				item = NewEC2KeyPairs(s.repos["EC2"].(*repo.EC2), s.app)
			case "Elasticache.Clusters":
				item = NewElasticacheClusters(s.repos["Elasticache"].(*repo.Elasticache), s.app)
			case "Elasticache.Events":
				item = NewElasticacheEvents(s.repos["Elasticache"].(*repo.Elasticache), s.app)
			case "Elasticache.Parameter Groups":
				item = NewElasticacheParameterGroups(s.repos["Elasticache"].(*repo.Elasticache), s.app)
			case "Elasticache.Reserved Nodes":
				item = NewElasticacheReservedCacheNodes(s.repos["Elasticache"].(*repo.Elasticache), s.app)
			case "Elasticache.Snapshots":
				item = NewElasticacheSnapshots(s.repos["Elasticache"].(*repo.Elasticache), s.app)
			case "ELB.Load Balancers":
				item = NewELBLoadBalancers(s.repos["ELB"].(*repo.ELB), s.app)
			case "ELB.Target Groups":
				item = NewELBTargetGroups(s.repos["ELB"].(*repo.ELB), s.app)
			case "IAM.Users":
				item = NewIAMUsers(s.repos["IAM"].(*repo.IAM), s.clients["IAM"].(*iam.Client), "", s.app)
			case "IAM.Roles":
				item = NewIAMRoles(s.repos["IAM"].(*repo.IAM), s.clients["IAM"].(*iam.Client), s.app)
			case "IAM.Groups":
				item = NewIAMGroups(s.repos["IAM"].(*repo.IAM), s.clients["IAM"].(*iam.Client), "", s.app)
			case "IAM.Managed Policies":
				item = NewIAMPolicies(s.repos["IAM"].(*repo.IAM), s.clients["IAM"].(*iam.Client), model.IAMIdentityTypeAll, "", s.app)
			case "KMS.Keys":
				item = NewKmsKeys(s.repos["KMS"].(*repo.KMS), s.app)
			case "KMS.Custom Key Stores":
				item = NewKmsCustomKeyStores(s.repos["KMS"].(*repo.KMS), s.app)
			case "Route 53.Hosted Zones":
				item = NewRoute53HostedZones(s.repos["Route 53"].(*repo.Route53), s.app)
			case "Route 53.Health Checks":
				item = NewRoute53HealthChecks(s.repos["Route 53"].(*repo.Route53), s.app)
			case "S3.Buckets":
				item = NewS3Buckets(s.repos["S3"].(*repo.S3), s.clients["S3"].(*s3.Client), s.app)
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
			return false
		}
		return true
	})
}

func (s Services) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s Services) Render() {
}
