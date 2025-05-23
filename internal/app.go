package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cw "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwLogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	ga "github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	msk "github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	sq "github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net/http"
	"strings"
)

type Application struct {
	app        *tview.Application
	pages      *tview.Pages
	header     *Header
	footer     *Footer
	components []Component
}

func NewApplication() *Application {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	westCfg := cfg.Copy()
	westCfg.Region = "us-west-2"

	app := tview.NewApplication()

	httpClient := &http.Client{}

	acmClient := acm.NewFromConfig(cfg)
	acmPCAClient := acmpca.NewFromConfig(cfg)
	cfClient := cf.NewFromConfig(cfg)
	cwClient := cw.NewFromConfig(cfg)
	cwLogsClient := cwLogs.NewFromConfig(cfg)
	ddbClient := ddb.NewFromConfig(cfg)
	ecClient := ec.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)
	ecsClient := ecs.NewFromConfig(cfg)
	eksClient := eks.NewFromConfig(cfg)
	elbClient := elb.NewFromConfig(cfg)
	gaClient := ga.NewFromConfig(westCfg)
	iamClient := iam.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)
	lambdaClient := lambda.NewFromConfig(cfg)
	mqClient := mq.NewFromConfig(cfg)
	mskClient := msk.NewFromConfig(cfg)
	rdsClient := rds.NewFromConfig(cfg)
	r53Client := r53.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)
	smClient := sm.NewFromConfig(cfg)
	sqClient := sq.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)
	ssmClient := ssm.NewFromConfig(cfg)
	stsClient := sts.NewFromConfig(cfg)

	a := &Application{}

	acmRepo := repo.NewACM(acmClient)
	acmPCARepo := repo.NewACMPCA(acmPCAClient)
	cfRepo := repo.NewCloudFront(cfClient)
	cwRepo := repo.NewCloudWatch(cwLogsClient)
	ddbRepo := repo.NewDynamoDB(ddbClient)
	ec2Repo := repo.NewEC2(ec2Client)
	ecsRepo := repo.NewECS(ecsClient)
	eksRepo := repo.NewEKS(eksClient)
	elbRepo := repo.NewELB(elbClient, httpClient)
	ecRepo := repo.NewElastiCache(ecClient, cwClient)
	gaRepo := repo.NewGlobalAccelerator(gaClient)
	iamRepo := repo.NewIAM(iamClient)
	kmsRepo := repo.NewKMS(kmsClient)
	lambdaRepo := repo.NewLambda(lambdaClient)
	mqRepo := repo.NewMQ(mqClient)
	mskRepo := repo.NewMSK(mskClient)
	rdsRepo := repo.NewRDS(rdsClient)
	r53Repo := repo.NewRoute53(r53Client)
	s3Repo := repo.NewS3(s3Client)
	snsRepo := repo.NewSNS(snsClient)
	smRepo := repo.NewSecretsManager(smClient)
	sqRepo := repo.NewServiceQuotas(sqClient)
	sqsRepo := repo.NewSQS(sqsClient)
	ssmRepo := repo.NewSSM(ssmClient)
	stsRepo := repo.NewSTS(stsClient)

	repos := map[string]interface{}{
		"ACM":                acmRepo,
		"ACM PCA":            acmPCARepo,
		"CloudFront":         cfRepo,
		"CloudWatch":         cwRepo,
		"DynamoDB":           ddbRepo,
		"EC2":                ec2Repo,
		"ECS":                ecsRepo,
		"EKS":                eksRepo,
		"ELB":                elbRepo,
		"ElastiCache":        ecRepo,
		"Global Accelerator": gaRepo,
		"IAM":                iamRepo,
		"MSK":                mskRepo,
		"KMS":                kmsRepo,
		"Lambda":             lambdaRepo,
		"MQ":                 mqRepo,
		"RDS":                rdsRepo,
		"Route 53":           r53Repo,
		"S3":                 s3Repo,
		"SNS":                snsRepo,
		"SQS":                sqsRepo,
		"STS":                stsRepo,
		"Secrets Manager":    smRepo,
		"SSM":                ssmRepo,
		"Service Quotas":     sqRepo,
	}

	services := NewServices(repos, a)
	pages := tview.NewPages()
	pages.SetBorder(true)

	header := NewHeader(stsRepo, iamRepo, a)
	footer := NewFooter(a)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(header, 4, 0, false) // header is 4 rows
	flex.AddItem(pages, 0, 1, true)   // main viewport is resizable
	flex.AddItem(footer, 1, 0, false) // footer is 1 row

	app.SetRoot(flex, true).SetFocus(pages)
	a.app = app
	a.pages = pages
	a.header = header
	a.footer = footer
	a.AddAndSwitch(services)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.Close()
			return nil
		}

		// pass down Enter keypress to the component
		if event.Key() == tcell.KeyEnter {
			return event
		}

		actions := a.GetActiveKeyActions()
		for _, action := range actions {
			if event.Name() == action.Key.Name() {
				action.Action()
				return nil
			}
		}
		return event
	})
	return a
}

func (a Application) refreshHandler() {
	_, primitive := a.pages.GetFrontPage()
	primitive.(Component).Render()
}

func (a Application) GetActiveKeyActions() []KeyAction {
	// TODO check that front page exists
	_, primitive := a.pages.GetFrontPage()
	// TODO avoid type coercion
	localActions := primitive.(Component).GetKeyActions()
	globalActions := []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModCtrl),
			Description: "Refresh",
			Action:      a.refreshHandler,
		},
	}
	return append(localActions, globalActions...)
}

func (a *Application) AddAndSwitch(v Component) {
	v.Render()
	// create a unique name for the tview pages element
	// TODO this hardcodes the index as part of the name to avoid collisions when similar views are chained together
	name := fmt.Sprintf("%v | %v | %v ", a.pages.GetPageCount(), v.GetService(), strings.Join(v.GetLabels(), " > "))
	a.components = append(a.components, v)
	a.pages.AddAndSwitchToPage(name, v, true)
	a.header.Render() // this has to happen after we update the pages view
	a.footer.Render()
	a.pages.SetTitle(fmt.Sprintf(" %v ", v.GetService()))
}

func (a *Application) Close() {
	// don't close if we're at the root page
	if a.pages.GetPageCount() == 1 {
		return
	}
	a.components = a.components[:len(a.components)-1]

	oldName, _ := a.pages.GetFrontPage()
	a.pages.RemovePage(oldName)
	// this assumes pages are retrieved in reverse order that they were added
	newName, _ := a.pages.GetFrontPage()
	a.pages.SwitchToPage(newName)
	a.pages.SetTitle(fmt.Sprintf(" %v ", a.components[len(a.components)-1].GetService()))
	a.header.Render()
	a.footer.Render()
}

func (a Application) Run() error {
	return a.app.Run()
}
