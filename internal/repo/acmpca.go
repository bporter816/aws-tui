package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/bporter816/aws-tui/internal/model"
)

type ACMPCA struct {
	acmPCAClient *acmpca.Client
}

func NewACMPCA(acmPCAClient *acmpca.Client) *ACMPCA {
	return &ACMPCA{
		acmPCAClient: acmPCAClient,
	}
}

func (a ACMPCA) ListCertificateAuthorities() ([]model.ACMPCACertificateAuthority, error) {
	pg := acmpca.NewListCertificateAuthoritiesPaginator(
		a.acmPCAClient,
		&acmpca.ListCertificateAuthoritiesInput{},
	)
	var certificateAuthorities []model.ACMPCACertificateAuthority
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ACMPCACertificateAuthority{}, err
		}
		for _, v := range out.CertificateAuthorities {
			certificateAuthorities = append(certificateAuthorities, model.ACMPCACertificateAuthority(v))
		}
	}
	return certificateAuthorities, nil
}

func (a ACMPCA) ListTags(certificateAuthorityArn string) (model.Tags, error) {
	pg := acmpca.NewListTagsPaginator(
		a.acmPCAClient,
		&acmpca.ListTagsInput{
			CertificateAuthorityArn: aws.String(certificateAuthorityArn),
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
