package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type S3Objects struct {
	*Tree
	s3Client *s3.Client
	bucket   string
	app      *Application
}

func NewS3Objects(s3Client *s3.Client, bucket string, app *Application) *S3Objects {
	root := tview.NewTreeNode("/")
	root.SetReference("")
	s := &S3Objects{
		Tree:     NewTree(root),
		s3Client: s3Client,
		bucket:   bucket,
		app:      app,
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s S3Objects) GetName() string {
	return fmt.Sprintf("S3 | %v | Objects", s.bucket)
}

func (s S3Objects) selectHandler(n *tview.TreeNode) {
	s.expandDir(n)
}

func (s S3Objects) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3Objects) expandDir(n *tview.TreeNode) {
	if strings.HasSuffix(n.GetText(), "/") {
		if len(n.GetChildren()) > 0 {
			n.SetExpanded(!n.IsExpanded())
			return
		}

		ref := n.GetReference().(string)
		pg := s3.NewListObjectsV2Paginator(
			s.s3Client,
			&s3.ListObjectsV2Input{
				Bucket:    aws.String(s.bucket),
				Delimiter: aws.String("/"),
				Prefix:    aws.String(ref),
			},
		)
		var prefixes []s3Types.CommonPrefix
		var objects []s3Types.Object
		for pg.HasMorePages() {
			out, err := pg.NextPage(context.TODO())
			if err != nil {
				panic(err)
			}
			prefixes = append(prefixes, out.CommonPrefixes...)
			objects = append(objects, out.Contents...)
		}
		for _, prefix := range prefixes {
			arr := strings.Split(*prefix.Prefix, "/")
			label := arr[len(arr)-2] + "/"
			c := tview.NewTreeNode(label)
			c.SetColor(tcell.ColorGreen)
			c.SetReference(ref + label)
			n.AddChild(c)
		}
		for _, object := range objects {
			key := *object.Key
			if strings.HasSuffix(key, "/") {
				continue
			}
			label := key[strings.LastIndex(key, "/")+1:]
			c := tview.NewTreeNode(label)
			c.SetReference(ref + label)
			n.AddChild(c)
		}
	} else {
		// TODO open file
	}
}

func (s S3Objects) Render() {
	s.expandDir(s.GetRoot())
}