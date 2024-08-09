package main

import (
	"strings"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type S3Objects struct {
	*ui.Tree
	view.S3
	repo   *repo.S3
	bucket string
	app    *Application
}

func NewS3Objects(repo *repo.S3, bucket string, app *Application) *S3Objects {
	root := tview.NewTreeNode(bucket + "/")
	root.SetReference("")
	s := &S3Objects{
		Tree:   ui.NewTree(root),
		repo:   repo,
		bucket: bucket,
		app:    app,
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s S3Objects) GetLabels() []string {
	return []string{s.bucket, "Objects"}
}

func (s S3Objects) selectHandler(n *tview.TreeNode) {
	s.expandDir(n)
}

func (s S3Objects) objectHandler() {
	if node := s.GetCurrentNode(); node != nil {
		key := node.GetReference().(string)
		if strings.HasSuffix(key, "/") {
			return
		}
		objectView := NewS3Object(s.repo, s.bucket, key, s.app)
		s.app.AddAndSwitch(objectView)
	}
}

func (s S3Objects) metadataHandler() {
	if node := s.GetCurrentNode(); node != nil {
		key := node.GetReference().(string)
		if strings.HasSuffix(key, "/") {
			return
		}
		metadataView := NewS3ObjectMetadata(s.repo, s.bucket, key, s.app)
		s.app.AddAndSwitch(metadataView)
	}
}

func (s S3Objects) tagsHandler() {
	if node := s.GetCurrentNode(); node != nil {
		key := node.GetReference().(string)
		if strings.HasSuffix(key, "/") {
			return
		}
		tagsView := NewTags(s.repo, s.GetService(), "object:"+s.bucket+":"+key, s.app)
		s.app.AddAndSwitch(tagsView)
	}
}

func (s S3Objects) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'v', tcell.ModNone),
			Description: "View Object",
			Action:      s.objectHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'm', tcell.ModNone),
			Description: "Metadata",
			Action:      s.metadataHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      s.tagsHandler,
		},
	}
}

func (s S3Objects) expandDir(n *tview.TreeNode) {
	if strings.HasSuffix(n.GetText(), "/") {
		if len(n.GetChildren()) > 0 {
			n.SetExpanded(!n.IsExpanded())
			return
		}

		ref := n.GetReference().(string)
		prefixes, objects, err := s.repo.ListObjects(s.bucket, ref)
		if err != nil {
			panic(err)
		}
		for _, prefix := range prefixes {
			arr := strings.Split(prefix, "/")
			label := arr[len(arr)-2] + "/"
			c := tview.NewTreeNode(label)
			c.SetColor(tcell.ColorGreen)
			c.SetReference(ref + label)
			n.AddChild(c)
		}
		for _, object := range objects {
			if strings.HasSuffix(object, "/") {
				continue
			}
			label := object[strings.LastIndex(object, "/")+1:]
			c := tview.NewTreeNode(label)
			c.SetReference(ref + label)
			n.AddChild(c)
		}
	}
}

func (s S3Objects) Render() {
	s.expandDir(s.GetRoot())
}
