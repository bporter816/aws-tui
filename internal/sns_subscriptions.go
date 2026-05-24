package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewSNSSubscriptions(repo *repo.SNS, topicArn string, app *Application) *TableView[model.SNSSubscription] {
	topicArnParsed, err := arn.Parse(topicArn)
	if err != nil {
		panic(err)
	}
	return NewTableView(
		"SNS",
		[]string{utils.GetResourceNameFromArn(topicArnParsed), "Subscriptions"},
		[]string{"ID", "PROTOCOL", "ENDPOINT", "STATUS"},
		func() ([]model.SNSSubscription, error) {
			return repo.ListSubscriptions(topicArn)
		},
		func(v model.SNSSubscription) []string {
			var id, protocol, status string
			if v.SubscriptionArn != nil {
				a, err := arn.Parse(*v.SubscriptionArn)
				if err != nil {
					panic(err)
				}
				id = utils.GetResourceNameFromArn(a)
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
