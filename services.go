package main

import (
	"fmt"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rivo/tview"
	"sort"
)

type Services struct {
	*Tree
	clients map[string]interface{}
	app     *Application
}

func NewServices(clients map[string]interface{}, app *Application) *Services {
	m := map[string][]string{
		"Cloudfront": []string{
			"Distributions",
		},
		"DynamoDB": []string{
			"Tables",
		},
		"Elasticache": []string{
			"Clusters",
			"Events",
			"Reserved Nodes",
		},
		"KMS": []string{
			"Keys",
		},
		"Route 53": []string{
			"Hosted Zones",
			"Health Checks",
		},
		"S3": []string{
			"Buckets",
		},
		"Secrets Manager": []string{
			"Secrets",
		},
	}
	root := tview.NewTreeNode("Services")
	s := &Services{
		Tree:    NewTree(root),
		clients: clients,
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

func (s Services) GetName() string {
	return "Services"
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

	s.root.Walk(func(node, parent *tview.TreeNode) bool {
		if node.GetText() == n.GetText() {
			view := fmt.Sprintf("%v.%v", parent.GetText(), node.GetText())
			var item Component
			switch view {
			case "Cloudfront.Distributions":
				item = NewCFDistributions(s.clients["Cloudfront"].(*cf.Client), s.app)
			case "DynamoDB.Tables":
				item = NewDynamoDBTables(s.clients["DynamoDB"].(*ddb.Client), s.app)
			case "Elasticache.Clusters":
				item = NewElasticacheClusters(s.clients["Elasticache"].(*ec.Client), s.app)
			case "Elasticache.Events":
				item = NewElasticacheEvents(s.clients["Elasticache"].(*ec.Client), s.app)
			case "Elasticache.Reserved Nodes":
				item = NewElasticacheReservedCacheNodes(s.clients["Elasticache"].(*ec.Client), s.app)
			case "KMS.Keys":
				item = NewKmsKeys(s.clients["KMS"].(*kms.Client), s.app)
			case "Route 53.Hosted Zones":
				item = NewRoute53HostedZones(s.clients["Route 53"].(*r53.Client), s.app)
			case "Route 53.Health Checks":
				item = NewRoute53HealthChecks(s.clients["Route 53"].(*r53.Client), s.app)
			case "S3.Buckets":
				item = NewS3Buckets(s.clients["S3"].(*s3.Client), s.app)
			case "Secrets Manager.Secrets":
				item = NewSMSecrets(s.clients["Secrets Manager"].(*sm.Client), s.app)
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
