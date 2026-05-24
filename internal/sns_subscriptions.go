package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewSNSSubscriptions(repo *repo.SNS, topicArn string, app *Application) *TableView[model.SNSSubscription] {
	topicLabel := topicArn
	if a, err := arn.Parse(topicArn); err == nil {
		topicLabel = utils.GetResourceNameFromArn(a)
	}
	return NewTableView(
		"SNS",
		[]string{topicLabel, "Subscriptions"},
		[]string{"ID", "PROTOCOL", "ENDPOINT", "STATUS"},
		func() ([]model.SNSSubscription, error) {
			return repo.ListSubscriptions(topicArn)
		},
		func(v model.SNSSubscription) []string {
			var id, protocol, status string
			if v.SubscriptionArn != nil {
				if a, err := arn.Parse(*v.SubscriptionArn); err == nil {
					id = utils.GetResourceNameFromArn(a)
				} else {
					id = *v.SubscriptionArn
				}
			}
			if v.Protocol != nil {
				protocol = utils.AutoCase(*v.Protocol)
			}
			if s, ok := v.Attributes["PendingConfirmation"]; ok {
				status = utils.BoolToString(s == "true", "Pending confirmation", "Confirmed")
			}
			return []string{
				id,
				protocol,
				utils.DerefString(v.Endpoint, ""),
				status,
			}
		},
		nil,
	)
}
