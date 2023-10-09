package model

import (
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type (
	EC2Instance          ec2Types.Instance
	EC2KeyPair           ec2Types.KeyPairInfo
	EC2SecurityGroup     ec2Types.SecurityGroup
	EC2SecurityGroupRule ec2Types.SecurityGroupRule
	EC2VPC               ec2Types.Vpc
	EC2Subnet            ec2Types.Subnet
	EC2AvailabilityZone  ec2Types.AvailabilityZone
)
