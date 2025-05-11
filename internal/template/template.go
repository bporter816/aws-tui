package template

import (
	"bytes"
	"crypto/x509"
	"github.com/bporter816/aws-tui/internal/utils"
	"math/big"
	"strings"
	"text/template"
	"time"
)

var (
	tmpl    *template.Template
	funcMap = template.FuncMap{
		"bigIntToBytes": func(b *big.Int) []byte {
			if b == nil {
				return []byte{}
			}
			return b.Bytes()
		},
		"chunk":        func(str string, length int) []string { return Chunk(str, length) },
		"formatSerial": func(bytes []byte) string { return FormatSerial(bytes) },
		"formatTime":   func(t time.Time) string { return t.Format(utils.DefaultTimeFormat) },
		"isRSA":        func(algo x509.PublicKeyAlgorithm) bool { return algo == x509.RSA },
		"join":         func(arr []string, sep string) string { return strings.Join(arr, sep) },
		"mul":          func(a, b int) int { return a * b },
	}
)

const (
	X509Certificate = "x509_certificate.tmpl"
)

func Init() {
	tmpl = template.Must(template.New("root").Funcs(funcMap).ParseGlob("internal/template/*.tmpl"))
}

func Render(templateName string, data any) (string, error) {
	var b bytes.Buffer
	err := tmpl.ExecuteTemplate(&b, templateName, data)
	return b.String(), err
}
