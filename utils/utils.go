package utils

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

var (
	titleCaser   = cases.Title(language.English)
	upperCaser   = cases.Upper(language.English)
	lowerCaser   = cases.Lower(language.English)
	re           = regexp.MustCompile("[-_]")
	sizePrefixes = []string{"KiB", "MiB", "GiB", "TiB", "PiB"}

	replacements = map[string]string{
		"ebs":   "EBS",
		"http":  "HTTP",
		"https": "HTTPS",
		"hvm":   "HVM",
		"iam":   "IAM",
		"ipv4":  "IPv4",
		"ipv6":  "IPv6",
		"json":  "JSON",
		"sms":   "SMS",
		"sqs":   "SQS",
	}

	abbreviations = map[string]string{
		"Milliseconds": "ms",
		"Seconds":      "sec",
		"Minutes":      "min",
		"Bytes":        "B",
		"Kilobytes":    "KB",
		"Megabytes":    "MB",
		"Gigabytes":    "GB",
	}
)

func TitleCase(str string) string {
	return titleCaser.String(str)
}

func UpperCase(str string) string {
	return upperCaser.String(str)
}

func LowerCase(str string) string {
	return lowerCaser.String(str)
}

func AutoCase(str string) string {
	strWithSpaces := re.ReplaceAllString(str, " ")
	words := strings.Split(strWithSpaces, " ")
	for i, w := range words {
		lower := LowerCase(w)
		if replacement, ok := replacements[lower]; ok {
			words[i] = replacement
		} else if i == 0 {
			words[i] = TitleCase(w)
		} else {
			words[i] = lower
		}
	}
	return strings.Join(words, " ")
}

func SimplifyFloat(value float64) string {
	if value == math.Trunc(value) {
		return strconv.Itoa(int(value))
	} else {
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
}

func AbbreviateUnit(unit string) string {
	if u, ok := abbreviations[unit]; ok {
		return u
	} else {
		return unit
	}
}

func BoolToString(b bool, y string, n string) string {
	if b {
		return y
	} else {
		return n
	}
}

func GetResourceNameFromArn(arn arn.ARN) string {
	// slash delimited resources
	parts := strings.Split(arn.Resource, "/")
	if len(parts) == 2 {
		return parts[len(parts)-1]
	} else if len(parts) > 2 {
		// path based resources
		idx := strings.Index(arn.Resource, "/")
		return arn.Resource[idx:]
	}

	// colon delimited resources
	parts = strings.Split(arn.Resource, ":")
	if len(parts) == 2 {
		return parts[len(parts)-1]
	}

	return arn.Resource
}

func ParseCertsFromPEM(data []byte) ([]*x509.Certificate, error) {
	var bytes []byte
	for {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}
		data = rest
		bytes = append(bytes, block.Bytes...)
	}
	return x509.ParseCertificates(bytes)
}

func FormatSize(size int64, precision int) string {
	// keep bytes as integers
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " B"
	}
	num := float64(size) / 1024.0
	for _, prefix := range sizePrefixes {
		if num < 1024.0 {
			return strconv.FormatFloat(num, 'f', precision, 64) + " " + prefix
		}
		num /= 1024.0
	}
	return "Too big"
}

func TruncateStrings(items []string, count int) string {
	diff := len(items) - count
	if diff <= 0 {
		return strings.Join(items, ", ")
	}
	return strings.Join(items[:count], ", ") + " + " + strconv.Itoa(diff) + " more"
}
