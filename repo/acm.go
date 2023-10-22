package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/bporter816/aws-tui/model"
)

type ACM struct {
	acmClient *acm.Client
}

func NewACM(acmClient *acm.Client) *ACM {
	return &ACM{
		acmClient: acmClient,
	}
}

func (a ACM) ListCertificates() ([]model.ACMCertificate, error) {
	pg := acm.NewListCertificatesPaginator(
		a.acmClient,
		&acm.ListCertificatesInput{
			Includes: &acmTypes.Filters{
				KeyTypes: acmTypes.KeyAlgorithmRsa2048.Values(),
			},
		},
	)
	var certificates []model.ACMCertificate
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ACMCertificate{}, err
		}
		for _, v := range out.CertificateSummaryList {
			certificates = append(certificates, model.ACMCertificate(v))
		}
	}
	return certificates, nil
}

func (a ACM) getCertificateDetails(certificateArn string) (*acm.GetCertificateOutput, error) {
	return a.acmClient.GetCertificate(
		context.TODO(),
		&acm.GetCertificateInput{
			CertificateArn: aws.String(certificateArn),
		},
	)
}

func (a ACM) GetCertificate(certificateArn string) (string, error) {
	out, err := a.getCertificateDetails(certificateArn)
	if err != nil {
		return "", err
	}
	if out.Certificate == nil {
		return "", nil
	}
	return *out.Certificate, nil
}

func (a ACM) ListTags(certificateArn arn.ARN) (model.Tags, error) {
	out, err := a.acmClient.ListTagsForCertificate(
		context.TODO(),
		&acm.ListTagsForCertificateInput{
			CertificateArn: aws.String(certificateArn.String()),
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
