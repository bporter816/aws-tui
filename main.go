package main

import (
	"context"
	"fmt"
	// "os/exec"
	// "sort"
	"strings"
	// "net/http"
	"github.com/aws/aws-sdk-go-v2/config"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	sq "github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/bporter816/aws-tui/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	app := tview.NewApplication()

	cfClient := cf.NewFromConfig(cfg)
	ddbClient := ddb.NewFromConfig(cfg)
	ecClient := ec.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)
	elbClient := elb.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)
	r53Client := r53.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	stsClient := sts.NewFromConfig(cfg)
	smClient := sm.NewFromConfig(cfg)
	sqClient := sq.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	a := &Application{}

	clients := map[string]interface{}{
		"DynamoDB":    ddbClient,
		"Elasticache": ecClient,
		"EC2":         ec2Client,
		"IAM":         iamClient,
		"KMS":         kmsClient,
		"Route 53":    r53Client,
		"S3":          s3Client,
		"STS":         stsClient,
	}

	repos := map[string]interface{}{
		"ELB":             repo.NewELB(elbClient),
		"Cloudfront":      repo.NewCloudfront(cfClient),
		"SQS":             repo.NewSQS(sqsClient),
		"Secrets Manager": repo.NewSecretsManager(smClient),
		"Service Quotas":  repo.NewServiceQuotas(sqClient),
	}

	services := NewServices(clients, repos, a)
	pages := tview.NewPages()
	pages.SetBorder(true)

	header := NewHeader(stsClient, iamClient, a)
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
		KeyAction{
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

func main() {
	// playing around with region selection panel, TODO move to dialog
	/*
		regionCmd := exec.Command("aws", "configure", "get", "region")
		regionOutput, err := regionCmd.Output()
		if err != nil {
			panic(err)
		}
		region := strings.TrimSpace(string(regionOutput))
		fmt.Printf("region: %v\n", region)

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic(err)
		}

		ec2Client := ec2.NewFromConfig(cfg)
		regions, err := ec2Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
		if err != nil {
			panic(err)
		}

		var regionsArr []string
		for _, r := range regions.Regions {
			if r.RegionName != nil {
				regionsArr = append(regionsArr, *r.RegionName)
			}
		}
		sort.Strings(regionsArr)

		l := tview.NewList()
		l.ShowSecondaryText(false)
		l.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'k' {
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			} else if event.Rune() == 'j' {
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			}
			return event
		})

		for i, v := range regionsArr {
			l.AddItem(v, "", 0, nil)
			if v == region {
				l.SetCurrentItem(i)
			}
		}
		l.SetOffset(0, 0)
	*/

	app := NewApplication()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
