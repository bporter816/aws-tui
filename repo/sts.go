package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/bporter816/aws-tui/model"
)

type STS struct {
	stsClient *sts.Client
}

func NewSTS(stsClient *sts.Client) *STS {
	return &STS{
		stsClient: stsClient,
	}
}

func (s STS) GetCallerIdentity() (model.STSCallerIdentity, error) {
	out, err := s.stsClient.GetCallerIdentity(
		context.TODO(),
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return model.STSCallerIdentity{}, err
	}
	return model.STSCallerIdentity(*out), nil
}
