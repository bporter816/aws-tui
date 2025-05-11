package repo

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bporter816/aws-tui/internal/model"
)

type S3 struct {
	s3Client *s3.Client
}

func NewS3(s3Client *s3.Client) *S3 {
	return &S3{
		s3Client: s3Client,
	}
}

func (s S3) ListBuckets() ([]model.S3Bucket, error) {
	out, err := s.s3Client.ListBuckets(
		context.TODO(),
		&s3.ListBucketsInput{},
	)
	if err != nil {
		return []model.S3Bucket{}, err
	}
	var buckets []model.S3Bucket
	for _, v := range out.Buckets {
		buckets = append(buckets, model.S3Bucket(v))
	}
	return buckets, nil
}

func (s S3) ListObjects(bucketName string, prefix string) ([]string, []string, error) {
	pg := s3.NewListObjectsV2Paginator(
		s.s3Client,
		&s3.ListObjectsV2Input{
			Bucket:    aws.String(bucketName),
			Delimiter: aws.String("/"),
			Prefix:    aws.String(prefix),
		},
	)
	var prefixes, objects []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []string{}, []string{}, err
		}
		for _, v := range out.CommonPrefixes {
			prefixes = append(prefixes, *v.Prefix)
		}
		for _, v := range out.Contents {
			objects = append(objects, *v.Key)
		}
	}
	return prefixes, objects, nil
}

func (s S3) GetBucketPolicy(bucketName string) (string, error) {
	out, err := s.s3Client.GetBucketPolicy(
		context.TODO(),
		&s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		},
	)
	if err != nil || out.Policy == nil {
		return "", err
	}
	return *out.Policy, nil
}

func (s S3) GetCORSRules(bucketName string) ([]model.S3CORSRule, error) {
	out, err := s.s3Client.GetBucketCors(
		context.TODO(),
		&s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		},
	)
	if err != nil {
		return []model.S3CORSRule{}, err
	}
	var corsRules []model.S3CORSRule
	for _, v := range out.CORSRules {
		corsRules = append(corsRules, model.S3CORSRule(v))
	}
	return corsRules, nil
}

func (s S3) listBucketTags(bucketName string) (model.Tags, error) {
	// TODO find where the panic occurs when there are no tags
	out, err := s.s3Client.GetBucketTagging(
		context.TODO(),
		&s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.TagSet {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}

func (s S3) GetObject(bucketName string, key string) ([]byte, error) {
	out, err := s.s3Client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return []byte{}, err
	}
	defer out.Body.Close()
	if out.ContentLength == nil {
		return []byte{}, errors.New("missing content length")
	}
	b, err := io.ReadAll(out.Body)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (s S3) GetObjectMetadata(bucketName string, key string) (model.Tags, error) {
	out, err := s.s3Client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	if out.ContentType != nil {
		tags = append(tags, model.Tag{Key: "Content-Type", Value: *out.ContentType})
	}
	for k, v := range out.Metadata {
		tags = append(tags, model.Tag{Key: k, Value: v})
	}
	return tags, nil
}

func (s S3) listObjectTags(bucketName string, key string) (model.Tags, error) {
	out, err := s.s3Client.GetObjectTagging(
		context.TODO(),
		&s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.TagSet {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}

func (s S3) ListTags(typeAndName string) (model.Tags, error) {
	parts := strings.Split(typeAndName, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return model.Tags{}, errors.New("must give type and name for s3 tags")
	}
	switch parts[0] {
	case "bucket":
		return s.listBucketTags(parts[1])
	case "object":
		return s.listObjectTags(parts[1], parts[2])
	default:
		return model.Tags{}, errors.New("must use bucket or object for s3 tags")
	}
}
