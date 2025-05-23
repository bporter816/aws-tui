package repo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/internal/model"
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

func (e EC2) GetPublicKey(keyPairId string) (string, error) {
	out, err := e.ec2Client.DescribeKeyPairs(
		context.TODO(),
		&ec2.DescribeKeyPairsInput{
			IncludePublicKey: aws.Bool(true),
			Filters: []ec2Types.Filter{
				{
					Name:   aws.String("key-pair-id"),
					Values: []string{keyPairId},
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	if len(out.KeyPairs) != 1 {
		return "", errors.New("should get exactly one key pair")
	}
	if out.KeyPairs[0].PublicKey == nil {
		return "", nil
	}
	return *out.KeyPairs[0].PublicKey, nil
}

func (e EC2) ListImages() ([]model.EC2Image, error) {
	pg := ec2.NewDescribeImagesPaginator(
		e.ec2Client,
		&ec2.DescribeImagesInput{
			Owners: []string{"self"},
		},
	)
	var images []model.EC2Image
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2Image{}, err
		}
		for _, v := range out.Images {
			images = append(images, model.EC2Image(v))
		}
	}
	return images, nil
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
				{
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

func (e EC2) ListSubnets(subnetIds []string) ([]model.EC2Subnet, error) {
	var pg *ec2.DescribeSubnetsPaginator
	if len(subnetIds) > 0 {
		pg = ec2.NewDescribeSubnetsPaginator(
			e.ec2Client,
			&ec2.DescribeSubnetsInput{
				SubnetIds: subnetIds,
			},
		)
	} else {
		pg = ec2.NewDescribeSubnetsPaginator(
			e.ec2Client,
			&ec2.DescribeSubnetsInput{},
		)
	}
	var subnets []model.EC2Subnet
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2Subnet{}, err
		}
		for _, v := range out.Subnets {
			subnets = append(subnets, model.EC2Subnet(v))
		}
	}
	return subnets, nil
}

func (e EC2) ListAvailabilityZones() ([]model.EC2AvailabilityZone, error) {
	out, err := e.ec2Client.DescribeAvailabilityZones(
		context.TODO(),
		&ec2.DescribeAvailabilityZonesInput{},
	)
	if err != nil {
		return []model.EC2AvailabilityZone{}, err
	}
	var availabilityZones []model.EC2AvailabilityZone
	for _, v := range out.AvailabilityZones {
		availabilityZones = append(availabilityZones, model.EC2AvailabilityZone(v))
	}
	return availabilityZones, nil
}

func (e EC2) ListReservedInstances(filters []ec2Types.Filter) ([]model.EC2ReservedInstance, error) {
	out, err := e.ec2Client.DescribeReservedInstances(
		context.TODO(),
		&ec2.DescribeReservedInstancesInput{
			Filters: filters,
		},
	)
	if err != nil {
		return []model.EC2ReservedInstance{}, err
	}
	var reservedInstances []model.EC2ReservedInstance
	for _, v := range out.ReservedInstances {
		reservedInstances = append(reservedInstances, model.EC2ReservedInstance(v))
	}
	return reservedInstances, nil
}

func (e EC2) ListInternetGateways() ([]model.EC2InternetGateway, error) {
	pg := ec2.NewDescribeInternetGatewaysPaginator(
		e.ec2Client,
		&ec2.DescribeInternetGatewaysInput{},
	)
	var internetGateways []model.EC2InternetGateway
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2InternetGateway{}, err
		}
		for _, v := range out.InternetGateways {
			internetGateways = append(internetGateways, model.EC2InternetGateway(v))
		}
	}
	return internetGateways, nil
}

func (e EC2) ListInternetGatewayAttachments(internetGatewayId string) ([]model.EC2InternetGatewayAttachment, error) {
	out, err := e.ec2Client.DescribeInternetGateways(
		context.TODO(),
		&ec2.DescribeInternetGatewaysInput{
			InternetGatewayIds: []string{internetGatewayId},
		},
	)
	if err != nil {
		return []model.EC2InternetGatewayAttachment{}, err
	}
	if len(out.InternetGateways) != 1 {
		return []model.EC2InternetGatewayAttachment{}, errors.New("should get exactly 1 internet gateway")
	}
	var attachments []model.EC2InternetGatewayAttachment
	for _, v := range out.InternetGateways[0].Attachments {
		attachments = append(attachments, model.EC2InternetGatewayAttachment(v))
	}
	return attachments, nil
}

func (e EC2) ListVolumes(filters []ec2Types.Filter) ([]model.EC2Volume, error) {
	pg := ec2.NewDescribeVolumesPaginator(
		e.ec2Client,
		&ec2.DescribeVolumesInput{
			Filters: filters,
		},
	)
	var volumes []model.EC2Volume
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EC2Volume{}, err
		}
		for _, v := range out.Volumes {
			volumes = append(volumes, model.EC2Volume(v))
		}
	}
	return volumes, nil
}

func (e EC2) ListTags(resourceId string) (model.Tags, error) {
	pg := ec2.NewDescribeTagsPaginator(
		e.ec2Client,
		&ec2.DescribeTagsInput{
			Filters: []ec2Types.Filter{
				{
					Name:   aws.String("resource-id"),
					Values: []string{resourceId},
				},
			},
		},
	)
	var tags model.Tags
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return model.Tags{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
		}
	}
	return tags, nil
}
