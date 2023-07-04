package main

import (
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gdamore/tcell/v2"
)

type Services struct {
	*Table
	clients map[string]interface{}
	app     *Application
}

func NewServices(clients map[string]interface{}, app *Application) *Services {
	s := &Services{
		Table: NewTable([]string{
			"SERVICE",
		}, 1, 0),
		clients: clients,
		app:     app,
	}
	return s
}

func (s Services) GetName() string {
	return "Services"
}

func (s Services) selectHandler() {
	service, err := s.GetColSelection("SERVICE")
	if err != nil {
		return
	}
	var item Component
	switch service {
	case "Cloudfront":
		item = NewCloudfrontDistributions(s.clients["Cloudfront"].(*cf.Client), s.app)
	case "DynamoDB":
		item = NewDynamoDBTables(s.clients["DynamoDB"].(*ddb.Client), s.app)
	case "Elasticache":
		item = NewElasticacheEvents(s.clients["Elasticache"].(*ec.Client), s.app)
	case "KMS":
		item = NewKmsKeys(s.clients["KMS"].(*kms.Client), s.app)
	case "Route 53":
		item = NewRoute53HostedZones(s.clients["Route 53"].(*r53.Client), s.app)
	case "S3":
		item = NewS3Buckets(s.clients["S3"].(*s3.Client), s.app)
	case "Secrets Manager":
		item = NewSecretsManagerSecrets(s.clients["Secrets Manager"].(*sm.Client), s.app)
	default:
		panic("unknown service")
	}
	s.app.AddAndSwitch(item)
}

func (s Services) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
			Description: "Select",
			Action:      s.selectHandler,
		},
	}
}

func (s Services) Render() {
	data := [][]string{
		[]string{"Cloudfront"},
		[]string{"DynamoDB"},
		[]string{"Elasticache"},
		[]string{"KMS"},
		[]string{"Route 53"},
		[]string{"S3"},
		[]string{"Secrets Manager"},
	}
	s.SetData(data)
}
