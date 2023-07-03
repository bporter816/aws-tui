package main

import (
	"context"
	"fmt"
	// "os/exec"
	// "sort"
	// "strings"
	// "net/http"
	"github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/ec2"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	// awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	// "github.com/aws/smithy-go/aws"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Application struct {
	app       *tview.Application
	pages     *tview.Pages
	pageNames []string
	header    *Header
}

func NewApplication() *Application {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	app := tview.NewApplication()

	stsClient := sts.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	r53Client := r53.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)
	smClient := sm.NewFromConfig(cfg)
	cfClient := cf.NewFromConfig(cfg)
	ddbClient := ddb.NewFromConfig(cfg)
	ecClient := ec.NewFromConfig(cfg)

	a := &Application{}

	clients := map[string]interface{}{
		"Cloudfront":      cfClient,
		"DynamoDB":        ddbClient,
		"Elasticache":     ecClient,
		"KMS":             kmsClient,
		"Route 53":        r53Client,
		"STS":             stsClient,
		"IAM":             iamClient,
		"Secrets Manager": smClient,
	}

	services := NewServices(clients, a)
	pages := tview.NewPages()
	pages.SetBorder(true)

	header := NewHeader(stsClient, iamClient, a)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(header, 4, 0, false)                          // header is 4 rows
	flex.AddItem(pages, 0, 1, true)                            // main viewport is resizable
	flex.AddItem(tview.NewTextView().SetText(""), 1, 0, false) // footer is 1 row

	app.SetRoot(flex, true).SetFocus(pages)
	a.app = app
	a.pages = pages
	a.header = header
	a.AddAndSwitch("services", services)
	a.pageNames = []string{"services"}
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape && len(a.pageNames) > 1 {
			a.Close()
			return nil
		}

		_, primitive := a.pages.GetFrontPage()
		if primitive != nil {
			actions := primitive.(Component).GetKeyActions()
			for _, action := range actions {
				if event.Name() == action.Key.Name() {
					action.Action()
					return nil
				}
			}
		}
		return event
	})
	return a
}

func (a Application) GetActiveKeyActions() []KeyAction {
	// TODO check that front page exists
	_, primitive := a.pages.GetFrontPage()
	// TODO avoid type coercion
	return primitive.(Component).GetKeyActions()
}

func (a *Application) AddAndSwitch(name string, v Component) {
	v.Render()
	a.pages.AddAndSwitchToPage(name, v, true)
	a.pageNames = append(a.pageNames, name)
	a.header.Render() // this has to happen after we update the pages view
	_, primitive := a.pages.GetFrontPage()
	a.pages.SetTitle(fmt.Sprintf(" %v ", primitive.(Component).GetName()))
}

func (a *Application) Close() {
	lastPageName := a.pageNames[len(a.pageNames)-1]
	a.pageNames = a.pageNames[:len(a.pageNames)-1]
	a.pages.RemovePage(lastPageName)
	a.pages.SwitchToPage(a.pageNames[len(a.pageNames)-1])
	a.header.Render()
	_, primitive := a.pages.GetFrontPage()
	a.pages.SetTitle(fmt.Sprintf(" %v ", primitive.(Component).GetName()))
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
