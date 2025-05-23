package repo

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/bporter816/aws-tui/internal/model"
)

type SecretsManager struct {
	smClient *sm.Client
}

func NewSecretsManager(smClient *sm.Client) *SecretsManager {
	return &SecretsManager{
		smClient: smClient,
	}
}

func (s SecretsManager) ListSecrets() ([]model.SecretsManagerSecret, error) {
	pg := sm.NewListSecretsPaginator(
		s.smClient,
		&sm.ListSecretsInput{
			IncludePlannedDeletion: aws.Bool(true),
		},
	)
	var secrets []model.SecretsManagerSecret
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.SecretsManagerSecret{}, err
		}
		for _, v := range out.SecretList {
			secrets = append(secrets, model.SecretsManagerSecret(v))
		}
	}
	return secrets, nil
}

func (s SecretsManager) GetSecretValue(secretName string) (string, error) {
	out, err := s.smClient.GetSecretValue(
		context.TODO(),
		&sm.GetSecretValueInput{
			SecretId: aws.String(secretName),
		},
	)
	if err != nil {
		return "", err
	}
	if out.SecretString != nil {
		return *out.SecretString, nil
	}
	if out.SecretBinary != nil {
		return base64.StdEncoding.EncodeToString(out.SecretBinary), nil
	}
	return "", nil
}

func (s SecretsManager) GetResourcePolicy(secretName string) (string, error) {
	out, err := s.smClient.GetResourcePolicy(
		context.TODO(),
		&sm.GetResourcePolicyInput{
			SecretId: aws.String(secretName),
		},
	)
	if err != nil {
		return "", err
	}
	var policy string
	if out.ResourcePolicy != nil {
		policy = *out.ResourcePolicy
	}
	return policy, nil
}

func (s SecretsManager) ListTags(secretName string) (model.Tags, error) {
	out, err := s.smClient.DescribeSecret(
		context.TODO(),
		&sm.DescribeSecretInput{
			SecretId: aws.String(secretName),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
