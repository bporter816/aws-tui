package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/bporter816/aws-tui/model"
)

type SNS struct {
	snsClient *sns.Client
}

func NewSNS(snsClient *sns.Client) *SNS {
	return &SNS{
		snsClient: snsClient,
	}
}

func (s SNS) getAttributes(topicArn string) (map[string]string, error) {
	out, err := s.snsClient.GetTopicAttributes(
		context.TODO(),
		&sns.GetTopicAttributesInput{
			TopicArn: aws.String(topicArn),
		},
	)
	if err != nil {
		return map[string]string{}, err
	}
	return out.Attributes, nil
}

func (s SNS) ListTopics() ([]model.SNSTopic, error) {
	pg := sns.NewListTopicsPaginator(
		s.snsClient,
		&sns.ListTopicsInput{},
	)
	var topics []model.SNSTopic
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.SNSTopic{}, err
		}
		for _, v := range out.Topics {
			if v.TopicArn == nil {
				continue
			}
			topic := model.SNSTopic{Arn: *v.TopicArn}
			if attrs, err := s.getAttributes(*v.TopicArn); err == nil {
				topic.Attributes = attrs
			}
			topics = append(topics, topic)
		}
	}
	return topics, nil
}

func (s SNS) ListTags(topicArn string) (model.Tags, error) {
	out, err := s.snsClient.ListTagsForResource(
		context.TODO(),
		&sns.ListTagsForResourceInput{
			ResourceArn: aws.String(topicArn),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
