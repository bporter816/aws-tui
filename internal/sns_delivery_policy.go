package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
)

type SNSDeliveryPolicy struct {
	*ui.Text
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

func (s SNSDeliveryPolicy) GetService() string { return "SNS" }

func (s SNSDeliveryPolicy) GetLabels() []string {
	label := s.topicArn
	if parsed, err := arn.Parse(s.topicArn); err == nil {
		label = utils.GetResourceNameFromArn(parsed)
	}
	return []string{label, "Delivery Policy"}
}

func (s SNSDeliveryPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSDeliveryPolicy) Render() error {
	policy, err := s.repo.GetDeliveryPolicy(s.topicArn)
	if err != nil {
		return err
	}
	s.SetText(policy)
	return nil
}
