package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/bporter816/aws-tui/model"
)

type Lambda struct {
	lambdaClient *lambda.Client
}

func NewLambda(lambdaClient *lambda.Client) *Lambda {
	return &Lambda{
		lambdaClient: lambdaClient,
	}
}

func (l Lambda) ListFunctions() ([]model.LambdaFunction, error) {
	pg := lambda.NewListFunctionsPaginator(
		l.lambdaClient,
		&lambda.ListFunctionsInput{},
	)
	var functions []model.LambdaFunction
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.LambdaFunction{}, err
		}
		for _, v := range out.Functions {
			functions = append(functions, model.LambdaFunction(v))
		}
	}
	return functions, nil
}

func (l Lambda) ListTags(functionArn string) (model.Tags, error) {
	out, err := l.lambdaClient.ListTags(
		context.TODO(),
		&lambda.ListTagsInput{
			Resource: aws.String(functionArn),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for k, v := range out.Tags {
		tags = append(tags, model.Tag{Key: k, Value: v})
	}
	return tags, nil
}
