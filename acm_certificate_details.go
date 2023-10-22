package main

import (
	"crypto/x509"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/template"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
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

func (a ACMCertificateDetails) GetService() string {
	return "ACM"
}

func (a ACMCertificateDetails) GetLabels() []string {
	arn, err := arn.Parse(a.certificateArn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), a.displayName}
}

func (a ACMCertificateDetails) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (a ACMCertificateDetails) Render() {
	cert, err := a.repo.GetCertificate(a.certificateArn)
	if err != nil {
		panic(err)
	}
	certs, err := utils.ParseCertsFromPEM([]byte(cert))
	if err != nil {
		panic(err)
	}
	text, err := template.Render(template.X509Certificate, struct {
		Metadata *x509.Certificate
		PEM      string
	}{
		Metadata: certs[0],
		PEM:      cert,
	})
	if err != nil {
		panic(err)
	}
	a.SetText(text)
}
