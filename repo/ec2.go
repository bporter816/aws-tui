package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/model"
)

type EC2 struct {
	ec2Client *ec2.Client
}

func NewEC2(ec2Client *ec2.Client) *EC2 {
	return &EC2{
		ec2Client: ec2Client,
	}
}

func (e EC2) ListInstances() ([]model.EC2Instance, error) {
	pg := ec2.NewDescribeInstancesPaginator(
		e.ec2Client,
		&ec2.DescribeInstancesInput{},
	)
	var instances []model.EC2Instance
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2Instance{}, err
		}
		for _, v := range out.Reservations {
			for _, vv := range v.Instances {
				instances = append(instances, model.EC2Instance(vv))
			}
		}
	}
	return instances, nil
}

func (e EC2) ListKeyPairs() ([]model.EC2KeyPair, error) {
	out, err := e.ec2Client.DescribeKeyPairs(
		context.TODO(),
		&ec2.DescribeKeyPairsInput{},
	)
	if err != nil {
		return []model.EC2KeyPair{}, err
	}
	var keyPairs []model.EC2KeyPair
	for _, v := range out.KeyPairs {
		keyPairs = append(keyPairs, model.EC2KeyPair(v))
	}
	return keyPairs, nil
}

func (e EC2) ListSecurityGroups() ([]model.EC2SecurityGroup, error) {
	pg := ec2.NewDescribeSecurityGroupsPaginator(
		e.ec2Client,
		&ec2.DescribeSecurityGroupsInput{},
	)
	var securityGroups []model.EC2SecurityGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2SecurityGroup{}, err
		}
		for _, v := range out.SecurityGroups {
			securityGroups = append(securityGroups, model.EC2SecurityGroup(v))
		}
	}
	return securityGroups, nil
}

func (e EC2) ListSecurityGroupRules(securityGroupId string) ([]model.EC2SecurityGroupRule, error) {
	pg := ec2.NewDescribeSecurityGroupRulesPaginator(
		e.ec2Client,
		&ec2.DescribeSecurityGroupRulesInput{
			Filters: []ec2Types.Filter{
				ec2Types.Filter{
					Name:   aws.String("group-id"),
					Values: []string{securityGroupId},
				},
			},
		},
	)
	var sgRules []model.EC2SecurityGroupRule
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2SecurityGroupRule{}, err
		}
		for _, v := range out.SecurityGroupRules {
			sgRules = append(sgRules, model.EC2SecurityGroupRule(v))
		}
	}
	return sgRules, nil
}

func (e EC2) ListVPCs() ([]model.EC2VPC, error) {
	pg := ec2.NewDescribeVpcsPaginator(
		e.ec2Client,
		&ec2.DescribeVpcsInput{},
	)
	var vpcs []model.EC2VPC
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2VPC{}, err
		}
		for _, v := range out.Vpcs {
			vpcs = append(vpcs, model.EC2VPC(v))
		}
	}
	return vpcs, nil
}
