package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/bporter816/aws-tui/model"
)

type MQ struct {
	mqClient *mq.Client
}

func NewMQ(mqClient *mq.Client) *MQ {
	return &MQ{
		mqClient: mqClient,
	}
}

func (m MQ) ListBrokers() ([]model.MQBroker, error) {
	pg := mq.NewListBrokersPaginator(
		m.mqClient,
		&mq.ListBrokersInput{},
	)
	var brokers []model.MQBroker
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.MQBroker{}, err
		}
		for _, v := range out.BrokerSummaries {
			brokers = append(brokers, model.MQBroker(v))
		}
	}
	return brokers, nil
}
