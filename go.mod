module github.com/bporter816/aws-tui

go 1.21

toolchain go1.22.5

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/aws/aws-sdk-go-v2 v1.32.5
	github.com/aws/aws-sdk-go-v2/config v1.28.5
	github.com/aws/aws-sdk-go-v2/service/acm v1.30.6
	github.com/aws/aws-sdk-go-v2/service/acmpca v1.37.7
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.42.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.43.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.43.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.37.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.192.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.52.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.52.1
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.43.3
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.42.0
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.29.6
	github.com/aws/aws-sdk-go-v2/service/iam v1.38.1
	github.com/aws/aws-sdk-go-v2/service/kms v1.37.6
	github.com/aws/aws-sdk-go-v2/service/lambda v1.67.0
	github.com/aws/aws-sdk-go-v2/service/mq v1.27.7
	github.com/aws/aws-sdk-go-v2/service/rds v1.91.0
	github.com/aws/aws-sdk-go-v2/service/route53 v1.46.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.67.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.34.6
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.25.6
	github.com/aws/aws-sdk-go-v2/service/sns v1.33.5
	github.com/aws/aws-sdk-go-v2/service/sqs v1.37.1
	github.com/aws/aws-sdk-go-v2/service/ssm v1.55.6
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.1
	github.com/gdamore/tcell/v2 v2.7.4
	github.com/rivo/tview v0.0.0-20241103174730-c76f7879f592
	golang.org/x/text v0.20.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.7 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.46 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.20 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/kafka v1.38.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.5 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/term v0.26.0 // indirect
)
