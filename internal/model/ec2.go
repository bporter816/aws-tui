package model

import (
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type (
	EC2Instance                  ec2Types.Instance
	EC2KeyPair                   ec2Types.KeyPairInfo
	EC2ReservedInstance          ec2Types.ReservedInstances
	EC2SecurityGroup             ec2Types.SecurityGroup
	EC2SecurityGroupRule         ec2Types.SecurityGroupRule
	EC2Image                     ec2Types.Image
	EC2VPC                       ec2Types.Vpc
	EC2Subnet                    ec2Types.Subnet
	EC2AvailabilityZone          ec2Types.AvailabilityZone
	EC2InternetGateway           ec2Types.InternetGateway
	EC2InternetGatewayAttachment ec2Types.InternetGatewayAttachment
	EC2Volume                    ec2Types.Volume
)
