package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
)

type SNSAccessControlPolicy struct {
	*ui.Text
	repo     *repo.SNS
	topicArn string
	app      *Application
}

func NewSNSAccessControlPolicy(repo *repo.SNS, topicArn string, app *Application) *SNSAccessControlPolicy {
	s := &SNSAccessControlPolicy{
		Text:     ui.NewText(true, "json"),
		repo:     repo,
		topicArn: topicArn,
		app:      app,
	}
	return s
}

func (s SNSAccessControlPolicy) GetService() string { return "SNS" }

func (s SNSAccessControlPolicy) GetLabels() []string {
	label := s.topicArn
	if parsed, err := arn.Parse(s.topicArn); err == nil {
		label = utils.GetResourceNameFromArn(parsed)
	}
	return []string{label, "Access Control Policy"}
}

func (s SNSAccessControlPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSAccessControlPolicy) Render(ctx context.Context) error {
	policy, err := s.repo.GetAccessControlPolicy(ctx, s.topicArn)
	if err != nil {
		return err
	}
	s.SetText(policy)
	return nil
}
