package main

import (
	"context"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CFFunctions struct {
	*Table
	cfClient *cf.Client
	app      *Application
}

func NewCFFunctions(cfClient *cf.Client, app *Application) *CFFunctions {
	c := &CFFunctions{
		Table: NewTable([]string{
			"NAME",
			"COMMENT",
			"STATUS",
			"STAGE",
			"CREATED",
			"MODIFIED",
		}, 1, 0),
		cfClient: cfClient,
		app:      app,
	}
	return c
}

func (c CFFunctions) GetService() string {
	return "Cloudfront"
}

func (c CFFunctions) GetLabels() []string {
	return []string{"Functions"}
}

func (c CFFunctions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFFunctions) Render() {
	// ListFunctions doesn't have a paginator
	var functions []cfTypes.FunctionSummary
	var marker *string
	for {
		out, err := c.cfClient.ListFunctions(
			context.TODO(),
			&cf.ListFunctionsInput{
				Marker: marker,
			},
		)
		if err != nil {
			panic(err)
		}
		functions = append(functions, out.FunctionList.Items...)
		marker = out.FunctionList.NextMarker
		if marker == nil {
			break
		}
	}

	caser := cases.Title(language.English)
	var data [][]string
	for _, v := range functions {
		var name, comment, status, stage, created, modified string
		if v.Name != nil {
			name = *v.Name
		}
		if v.FunctionConfig != nil && v.FunctionConfig.Comment != nil {
			comment = *v.FunctionConfig.Comment
		}
		if v.Status != nil {
			status = caser.String(*v.Status)
		}
		if v.FunctionMetadata != nil {
			stage = caser.String(string(v.FunctionMetadata.Stage))
			if v.FunctionMetadata.CreatedTime != nil {
				created = v.FunctionMetadata.CreatedTime.Format("2006-01-02 15:04:05")
			}
			if v.FunctionMetadata.LastModifiedTime != nil {
				modified = v.FunctionMetadata.LastModifiedTime.Format("2006-01-02 15:04:05")
			}
		}
		data = append(data, []string{
			name,
			comment,
			status,
			stage,
			created,
			modified,
		})
	}
	c.SetData(data)
}
