package internal

import (
	"fmt"
	"sort"

	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/rivo/tview"
)

type Services struct {
	*ui.Tree
	repos Repos
	app   *Application
}

func NewServices(repos Repos, app *Application) *Services {
	root := tview.NewTreeNode("Services")
	s := &Services{
		Tree:  ui.NewTree(root),
		repos: repos,
		app:   app,
	}
	m := serviceTree()
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		n := tview.NewTreeNode(k)
		root.AddChild(n)
		for _, view := range m[k] {
			leaf := tview.NewTreeNode(view)
			leaf.SetReference(fmt.Sprintf("%v.%v", k, view))
			n.AddChild(leaf)
		}
		n.CollapseAll()
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s Services) GetService() string {
	return "Services"
}

func (s Services) GetLabels() []string {
	return []string{}
}

func (s Services) selectHandler(n *tview.TreeNode) {
	if n.GetLevel() < 2 {
		if n.IsExpanded() {
			n.Collapse()
		} else {
			n.Expand()
			if len(n.GetChildren()) > 0 {
				s.SetCurrentNode(n.GetChildren()[0])
			}
		}
		return
	}

	key := n.GetReference().(string)
	f, ok := registry[key]
	if !ok {
		return
	}
	s.app.AddAndSwitch(f(s.repos, s.app))
}

func (s Services) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s Services) Render() {
}
