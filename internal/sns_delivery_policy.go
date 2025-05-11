package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type SNSDeliveryPolicy struct {
	*ui.Text
	view.SNS
	repo     *repo.SNS
	topicArn string
	app      *Application
}

func NewSNSDeliveryPolicy(repo *repo.SNS, topicArn string, app *Application) *SNSDeliveryPolicy {
	s := &SNSDeliveryPolicy{
		Text:     ui.NewText(true, "json"),
		repo:     repo,
		topicArn: topicArn,
		app:      app,
	}
	return s
}

func (s SNSDeliveryPolicy) GetLabels() []string {
	arn, err := arn.Parse(s.topicArn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Delivery Policy"}
}

func (s SNSDeliveryPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSDeliveryPolicy) Render() {
	policy, err := s.repo.GetDeliveryPolicy(s.topicArn)
	if err != nil {
		panic(err)
	}
	s.SetText(policy)
}
