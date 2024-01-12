# aws-tui

A k9s-inspired terminal UI for AWS services. Supports:

* ACM
* CloudFront
* CloudWatch
* DynamoDB
* EBS
* EC2
* EKS
* ELB
* ElastiCache
* IAM
* Key Management Service
* Lambda
* Route 53
* S3
* SNS
* SQS
* Secrets Manager
* Service Quotas
* VPC

Note: There is a dependency on the AWS CLI for operations that are not easily supported by the AWS Go SDK:
* determining the default region for the current profile
* listing profiles (not used yet)
