package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/model"
)

type SQS struct {
	sqsClient *sqs.Client
}

func NewSQS(sqsClient *sqs.Client) *SQS {
	return &SQS{
		sqsClient: sqsClient,
	}
}

func (s SQS) getAttributes(queueUrl string, attributeNames []sqsTypes.QueueAttributeName) (map[string]string, error) {
	out, err := s.sqsClient.GetQueueAttributes(
		context.TODO(),
		&sqs.GetQueueAttributesInput{
			QueueUrl: aws.String(queueUrl),
			AttributeNames: attributeNames,
		},
	)
	if err != nil {
		return map[string]string{}, err
	}
	return out.Attributes, nil
}

func (s SQS) ListQueues() ([]model.SQSQueue, error) {
	var queues []model.SQSQueue
	pg := sqs.NewListQueuesPaginator(
		s.sqsClient,
		&sqs.ListQueuesInput{},
	)
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.SQSQueue{}, err
		}
		for _, v := range out.QueueUrls {
			attrs, err := s.getAttributes(v, []sqsTypes.QueueAttributeName{sqsTypes.QueueAttributeNameAll})
			if err != nil {
				continue
			}
			var queue model.SQSQueue
			// get the ARN and pull the name out
			if queueArn, ok := attrs[string(sqsTypes.QueueAttributeNameQueueArn)]; ok {
				arn, err := arn.Parse(queueArn)
				if err == nil {
					queue.Name = arn.Resource
				}
			}
			queue.QueueUrl = v
			if isFifo, ok := attrs[string(sqsTypes.QueueAttributeNameFifoQueue)]; ok && isFifo == "true" {
				queue.IsFifo = true
			}
			queues = append(queues, queue)
		}
	}
	return queues, nil
}

func (s SQS) GetPolicy(queueUrl string) (string, error) {
	attrs, err := s.getAttributes(queueUrl, []sqsTypes.QueueAttributeName{sqsTypes.QueueAttributeNamePolicy})
	if err != nil {
		return "", err
	}
	if policy, ok := attrs[string(sqsTypes.QueueAttributeNamePolicy)]; ok {
		return policy, nil
	}
	return "", nil
}

func (s SQS) ListTags(queueUrl string) (model.Tags, error) {
	out, err := s.sqsClient.ListQueueTags(
		context.TODO(),
		&sqs.ListQueueTagsInput{
			QueueUrl: aws.String(queueUrl),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for k, v := range out.Tags {
		tags = append(tags, model.Tag{Key: k, Value: v})
	}
	return tags, nil
}
