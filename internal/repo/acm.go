package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/bporter816/aws-tui/internal/model"
)

type ACM struct {
	acmClient *acm.Client
}

func NewACM(acmClient *acm.Client) *ACM {
	return &ACM{
		acmClient: acmClient,
	}
}

func (a ACM) ListCertificates(ctx context.Context) ([]model.ACMCertificate, error) {
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
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []model.ACMCertificate{}, err
		}
		for _, v := range out.CertificateSummaryList {
			certificates = append(certificates, model.ACMCertificate(v))
		}
	}
	return certificates, nil
}

func (a ACM) getCertificateDetails(ctx context.Context, certificateArn string) (*acm.GetCertificateOutput, error) {
	return a.acmClient.GetCertificate(
		ctx,
		&acm.GetCertificateInput{
			CertificateArn: aws.String(certificateArn),
		},
	)
}

func (a ACM) GetCertificate(ctx context.Context, certificateArn string) (string, error) {
	out, err := a.getCertificateDetails(ctx, certificateArn)
	if err != nil {
		return "", err
	}
	if out.Certificate == nil {
		return "", nil
	}
	return *out.Certificate, nil
}

func (a ACM) ListTags(ctx context.Context, resourceId string) (model.Tags, error) {
	out, err := a.acmClient.ListTagsForCertificate(
		ctx,
		&acm.ListTagsForCertificateInput{
			CertificateArn: aws.String(resourceId),
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
