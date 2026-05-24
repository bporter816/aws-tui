package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	aws "github.com/aws/aws-sdk-go-v2/aws"
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
)

type Application struct {
	app        *tview.Application
	pages      *tview.Pages
	header     *Header
	footer     *Footer
	components []Component
	region     string
	profile    string
}

func buildRepos(cfg aws.Config) (Repos, *repo.STS) {
	// Global Accelerator requires us-west-2 regardless of selected region.
	westCfg := cfg.Copy()
	westCfg.Region = "us-west-2"

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

	repos := Repos{
		ACM:               repo.NewACM(acmClient),
		ACMPCA:            repo.NewACMPCA(acmPCAClient),
		CloudFront:        repo.NewCloudFront(cfClient),
		CloudWatch:        repo.NewCloudWatch(cwLogsClient),
		DynamoDB:          repo.NewDynamoDB(ddbClient),
		EC2:               repo.NewEC2(ec2Client),
		ECS:               repo.NewECS(ecsClient),
		EKS:               repo.NewEKS(eksClient),
		ELB:               repo.NewELB(elbClient, httpClient),
		ElastiCache:       repo.NewElastiCache(ecClient, cwClient),
		GlobalAccelerator: repo.NewGlobalAccelerator(gaClient),
		IAM:               repo.NewIAM(iamClient),
		KMS:               repo.NewKMS(kmsClient),
		Lambda:            repo.NewLambda(lambdaClient),
		MQ:                repo.NewMQ(mqClient),
		MSK:               repo.NewMSK(mskClient),
		RDS:               repo.NewRDS(rdsClient),
		Route53:           repo.NewRoute53(r53Client),
		S3:                repo.NewS3(s3Client),
		SecretsManager:    repo.NewSecretsManager(smClient),
		ServiceQuotas:     repo.NewServiceQuotas(sqClient),
		SNS:               repo.NewSNS(snsClient),
		SQS:               repo.NewSQS(sqsClient),
		SSM:               repo.NewSSM(ssmClient),
	}
	stsRepo := repo.NewSTS(stsClient)
	return repos, stsRepo
}

func NewApplication() *Application {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	repos, stsRepo := buildRepos(cfg)

	a := &Application{
		region: cfg.Region,
	}

	tviewApp := tview.NewApplication()
	pages := tview.NewPages()
	pages.SetBorder(true)

	header := NewHeader(stsRepo, repos.IAM, a)
	footer := NewFooter(a)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(header, 4, 0, false) // header is 4 rows
	flex.AddItem(pages, 0, 1, true)   // main viewport is resizable
	flex.AddItem(footer, 1, 0, false) // footer is 1 row

	tviewApp.SetRoot(flex, true).SetFocus(pages)
	a.app = tviewApp
	a.pages = pages
	a.header = header
	a.footer = footer
	a.AddAndSwitch(NewServices(repos, a))

	tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		frontName, _ := a.pages.GetFrontPage()

		// dismiss overlays with Escape without disturbing the navigation stack
		if frontName == "picker" || frontName == "error" {
			if event.Key() == tcell.KeyEscape {
				a.pages.RemovePage(frontName)
				a.app.SetFocus(a.pages)
				return nil
			}
			return event
		}

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

func (a *Application) GetRegion() string {
	return a.region
}

func (a *Application) rebuild(profile, region string) {
	opts := []func(*config.LoadOptions) error{}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		a.ShowError(err)
		return
	}

	a.profile = profile
	a.region = cfg.Region

	repos, stsRepo := buildRepos(cfg)
	a.header.UpdateRepos(stsRepo, repos.IAM)

	// Walk the current navigation stack backward to find the deepest component
	// that is a generic list view (i.e. its key exists in the registry). Detail
	// views for specific resources have a resource ID as their first label, so
	// they won't match. We stop at the first match and recreate from there.
	var resumeKey string
	for i := len(a.components) - 1; i >= 0; i-- {
		labels := a.components[i].GetLabels()
		if len(labels) == 0 {
			continue
		}
		key := a.components[i].GetService() + "." + labels[0]
		if _, ok := registry[key]; ok {
			resumeKey = key
			break
		}
	}

	for a.pages.GetPageCount() > 0 {
		name, _ := a.pages.GetFrontPage()
		if name == "" {
			break
		}
		a.pages.RemovePage(name)
	}
	a.components = nil

	a.AddAndSwitch(NewServices(repos, a))
	if resumeKey != "" {
		a.AddAndSwitch(registry[resumeKey](repos, a))
	}
}

func (a *Application) ShowError(err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorDefault).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			a.pages.RemovePage("error")
			a.app.SetFocus(a.pages)
		})
	modal.Box.SetBackgroundColor(tcell.ColorDefault)
	a.pages.AddPage("error", modal, false, true)
	a.app.SetFocus(modal)
}

func (a *Application) refreshHandler() {
	_, primitive := a.pages.GetFrontPage()
	if err := primitive.(Component).Render(); err != nil {
		a.ShowError(err)
	}
}

func (a *Application) GetActiveKeyActions() []KeyAction {
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
		{
			Key:         tcell.NewEventKey(tcell.KeyCtrlG, 0, tcell.ModCtrl),
			Description: "Switch Region",
			Action:      a.showRegionPicker,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyCtrlP, 0, tcell.ModCtrl),
			Description: "Switch Profile",
			Action:      a.showProfilePicker,
		},
	}
	return append(localActions, globalActions...)
}

func (a *Application) AddAndSwitch(v Component) {
	if err := v.Render(); err != nil {
		a.ShowError(err)
		return
	}
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
