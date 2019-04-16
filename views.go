package main

import (
	"io"
	. "yumaikas/werewolf/templates"
)

func render(w io.Writer, template func(Context)) {
	die(RenderWithTargetAndTheme(w, "AQUA", template))
}

func HomePageView(w io.Writer, nodes []OutlineNodeDB) {
	// TODO: Clean this up to not use a special node ID. Maybe
	// have a rootlist?
	render(w, BasePage("Outliner Home",
		renderNodesInOutlineOrder(NodesToTree(nodes))))
}

// New plan: Turn the list of nodes into a tree, and recursive render that sucker,
// using a BFS traversal, because this fiddly stuff is escaping my mind too much.
func renderNodesInOutlineOrder(nodes []*OutlineTree) func(Context) {
	return func(ctx Context) {
		for _, n := range nodes {
			renderNodeTree(*n)(ctx)
		}
	}
}

func renderNodeTree(n OutlineTree) func(Context) {
	content := StringOr(n.Self.Content, "")
	if len(n.Children) > 0 {
		return Div(Atr, Str(content), func(ctx Context) {
			for _, c := range n.Children {
				renderNodeTree(*c)(ctx)
			}
		})
	} else {
		return Div(Atr, Str(content))
	}
}
