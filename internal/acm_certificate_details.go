package internal

import (
	"context"
	"crypto/x509"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/template"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
)


type ACMCertificateDetails struct {
	*ui.Text
	repo           *repo.ACM
	displayName    string
	certificateArn string
	app            *Application
}

func NewACMCertificateDetails(repo *repo.ACM, displayName string, certificateArn string, app *Application) *ACMCertificateDetails {
	a := &ACMCertificateDetails{
		Text:           ui.NewText(false, "text"),
		repo:           repo,
		displayName:    displayName,
		certificateArn: certificateArn,
		app:            app,
	}
	return a
}

func (a ACMCertificateDetails) GetService() string { return "ACM" }

func (a ACMCertificateDetails) GetLabels() []string {
	label := a.certificateArn
	if parsed, err := arn.Parse(a.certificateArn); err == nil {
		label = utils.GetResourceNameFromArn(parsed)
	}
	return []string{label, a.displayName}
}

func (a ACMCertificateDetails) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (a ACMCertificateDetails) Render(ctx context.Context) error {
	cert, err := a.repo.GetCertificate(ctx, a.certificateArn)
	if err != nil {
		return err
	}
	certs, err := utils.ParseCertsFromPEM([]byte(cert))
	if err != nil {
		return err
	}
	text, err := template.Render(template.X509Certificate, struct {
		Metadata *x509.Certificate
		PEM      string
	}{
		Metadata: certs[0],
		PEM:      cert,
	})
	if err != nil {
		return err
	}
	a.SetText(text)
	return nil
}
