package repo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/bporter816/aws-tui/internal/model"
)

type ELB struct {
	elbClient  *elb.Client
	httpClient *http.Client
}

func NewELB(elbClient *elb.Client, httpClient *http.Client) *ELB {
	return &ELB{
		elbClient:  elbClient,
		httpClient: httpClient,
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

func (e ELB) ListTags(resourceArn string) (model.Tags, error) {
	out, err := e.elbClient.DescribeTags(
		context.TODO(),
		&elb.DescribeTagsInput{
			ResourceArns: []string{resourceArn},
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	if len(out.TagDescriptions) != 1 {
		return model.Tags{}, errors.New("should get exactly 1 tag description")
	}
	var tags model.Tags
	for _, v := range out.TagDescriptions[0].Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}

func (e ELB) ListTrustStores() ([]model.ELBTrustStore, error) {
	pg := elb.NewDescribeTrustStoresPaginator(
		e.elbClient,
		&elb.DescribeTrustStoresInput{},
	)
	var trustStores []model.ELBTrustStore
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ELBTrustStore{}, err
		}
		for _, v := range out.TrustStores {
			trustStores = append(trustStores, model.ELBTrustStore(v))
		}
	}
	return trustStores, nil
}

func (e ELB) ListTrustStoreAssociations(a arn.ARN) ([]model.ELBTrustStoreAssociation, error) {
	pg := elb.NewDescribeTrustStoreAssociationsPaginator(
		e.elbClient,
		&elb.DescribeTrustStoreAssociationsInput{
			TrustStoreArn: aws.String(a.String()),
		},
	)
	var associations []model.ELBTrustStoreAssociation
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ELBTrustStoreAssociation{}, err
		}
		for _, v := range out.TrustStoreAssociations {
			associations = append(associations, model.ELBTrustStoreAssociation(v))
		}
	}
	return associations, nil
}

func (e ELB) GetTrustStoreCACertificatesBundle(trustStoreArn string) (string, error) {
	out, err := e.elbClient.GetTrustStoreCaCertificatesBundle(
		context.TODO(),
		&elb.GetTrustStoreCaCertificatesBundleInput{
			TrustStoreArn: aws.String(trustStoreArn),
		},
	)
	if err != nil {
		return "", err
	}
	if out.Location == nil {
		return "", errors.New("empty s3 location for trust store")
	}
	resp, err := e.httpClient.Get(*out.Location)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
