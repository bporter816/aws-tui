package ui

import (
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.TreeView
	Root *tview.TreeNode
}

func NewTree(root *tview.TreeNode) *Tree {
	tv := tview.NewTreeView()
	tv.SetGraphics(true)
	tv.SetRoot(root)
	tv.SetCurrentNode(root)
	t := &Tree{
		TreeView: tv,
		Root:     root,
	}
	return t
}
