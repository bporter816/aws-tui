package repo

import (
	"context"
	"errors"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/bporter816/aws-tui/model"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type ELB struct {
	elbClient *elb.Client
}

func NewELB(elbClient *elb.Client) *ELB {
	return &ELB{
		elbClient: elbClient,
	}
}

func (e ELB) ListLoadBalancers() ([]model.ELBLoadBalancer, error) {
	pg := elb.NewDescribeLoadBalancersPaginator(
		e.elbClient,
		&elb.DescribeLoadBalancersInput{},
	)
	var loadBalancers []model.ELBLoadBalancer
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ELBLoadBalancer{}, err
		}
		for _, v := range out.LoadBalancers {
			loadBalancers = append(loadBalancers, model.ELBLoadBalancer(v))
		}
	}
	return loadBalancers, nil
}

func (e ELB) ListListeners(loadBalancerArn string) ([]model.ELBListener, error) {
	pg := elb.NewDescribeListenersPaginator(
		e.elbClient,
		&elb.DescribeListenersInput{
			LoadBalancerArn: aws.String(loadBalancerArn),
		},
	)
	var listeners []model.ELBListener
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ELBListener{}, err
		}
		for _, v := range out.Listeners {
			m := model.ELBListener{Listener: v}
			if v.ListenerArn != nil {
				if rules, err := e.ListListenerRules(*v.ListenerArn); err != nil {
					m.Rules = len(rules)
				}
			}
			listeners = append(listeners, m)
		}
	}
	return listeners, nil
}

func (e ELB) ListListenerRules(listenerArn string) ([]model.ELBListenerRule, error) {
	var marker *string
	var rules []model.ELBListenerRule
	for {
		out, err := e.elbClient.DescribeRules(
			context.TODO(),
			&elb.DescribeRulesInput{
				ListenerArn: aws.String(listenerArn),
				Marker:      marker,
			},
		)
		if err != nil {
			return []model.ELBListenerRule{}, err
		}
		for _, v := range out.Rules {
			rules = append(rules, model.ELBListenerRule(v))
		}
		marker = out.NextMarker
		if marker == nil {
			break
		}
	}
	return rules, nil
}

func (e ELB) ListTargetGroups() ([]model.ELBTargetGroup, error) {
	pg := elb.NewDescribeTargetGroupsPaginator(
		e.elbClient,
		&elb.DescribeTargetGroupsInput{},
	)
	var targetGroups []model.ELBTargetGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ELBTargetGroup{}, err
		}
		for _, v := range out.TargetGroups {
			targetGroups = append(targetGroups, model.ELBTargetGroup(v))
		}
	}
	return targetGroups, nil
}

func (e ELB) ListTags(resourceArn string) ([]model.Tag, error) {
	out, err := e.elbClient.DescribeTags(
		context.TODO(),
		&elb.DescribeTagsInput{
			ResourceArns: []string{resourceArn},
		},
	)
	if err != nil {
		return []model.Tag{}, err
	}
	if len(out.TagDescriptions) != 1 {
		return []model.Tag{}, errors.New("should get exactly 1 tag description")
	}
	var tags []model.Tag
	for _, v := range out.TagDescriptions[0].Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
