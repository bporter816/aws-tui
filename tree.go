package main

import (
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.TreeView
	root *tview.TreeNode
}

func NewTree(root *tview.TreeNode) *Tree {
	tv := tview.NewTreeView()
	tv.SetGraphics(true)
	tv.SetRoot(root)
	tv.SetCurrentNode(root)
	t := &Tree{
		TreeView: tv,
		root:     root,
	}
	return t
}
