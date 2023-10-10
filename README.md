# aws-tui

A k9s-inspired terminal UI for AWS services. Supports:

* Cloudfront
* DynamoDB
* EC2 (VPCs and Security Groups)
* ELB
* ElastiCache
* IAM
* Key Management Service
* Route 53
* S3
* Secrets Manager
* Service Quotas

Note: There is a dependency on the AWS CLI for operations that are not easily supported by the AWS Go SDK:
* determining the default region for the current profile
* listing profiles (not used yet)
