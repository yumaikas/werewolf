package main

import (
	"fmt"
	"io"
	"io/ioutil"
	. "werewolf/templates"
)

func render(w io.Writer, template func(Context)) {
	die(RenderWithTargetAndTheme(w, "AQUA", template))
}

// Use to make editing files and reloading things easier
func WIPFile(path string) string {
	content, err := ioutil.ReadFile(path)
	die(err)
	return string(content)
}

var outlineStyles = `
.outline-node-top {
	margin-left: 5px;
}
.outline-node-inner {
	margin-left: 20px;
}
.hidden {
	display: none;
}
`

func HomePageView(w io.Writer, nodes []OutlineNodeDB) {
	// This takes a list of outline roots and renders them
	render(w, BasePage("Outliner Home",
		Style(outlineStyles),
		renderNodesInOutlineOrder(NodesToTree(nodes)),
		JS(WIPFile("static/reqwest.js")),
		outlineScripts()))
}

func outlineScripts() func(Context) {
	return JS(WIPFile("WIPOutlineScripts.js"))
}

// Render using using a BFS traversal, because this fiddly stuff is escaping my mind too much.
func renderNodesInOutlineOrder(nodes []*OutlineTree) func(Context) {
	return func(ctx Context) {
		// Get the top node, and then emit a link to it, calling it "UP"
		if len(nodes) > 0 {
			top := *nodes[0]
			Div(Atr.Id("top-bar"),
				// Hide this div by default
				Div(Atr.Id("message-bar").Class("hidden"), Str("If you see this, it's a bug")),
				Div(Atr.Id("command-bar"),
					Str("Commands: "),
					A(Atr.Href(parentLink(top)), RawStr("Up")),
					A(Atr.UnsafeHref("javascript:werewolf.expandAll();"), RawStr("Expand all")),
					A(Atr.UnsafeHref("javascript:werewolf.collapseAll();"), RawStr("Collapse all")),
					A(Atr.UnsafeHref("javascript:werewolf.beginCreateNode();"), RawStr("Create Node")),
				),
				Hr(),
			)(ctx)

			for _, n := range nodes {
				renderNodeTree(*n)(ctx)
			}
		}
	}
}

func nodeRecur(n OutlineTree) func(Context) {
	return func(ctx Context) {
		for _, c := range n.Children {
			renderNodeTree(*c)(ctx)
		}
	}
}

func parentLink(n OutlineTree) string {
	// If we can go up
	if n.Self.ParentId.Valid {
		return fmt.Sprint("/node/", n.Self.ParentId.Int64, "/page")
	}
	// Otherwise, keep it here.
	return pageLink(n)
}

func pageLink(n OutlineTree) string {
	return fmt.Sprint("/node/", n.Self.Id, "/page")
}

// Emit the various links for a node
func nodeInfo(n OutlineTree) func(Context) {
	return A(Atr.Href(pageLink(n)), Str("(Zoom)"))
}

func renderNodeTree(n OutlineTree) func(Context) {
	content := StringOr(n.Self.Content, "")
	id := n.Self.Id
	var class = "outline-node"
	if n.Self.RelativeDepth == 0 {
		class += " outline-node-top"
	} else {
		class += " outline-node-inner"
	}
	var inner = nodeRecur(n)
	if len(n.Children) <= 0 {
		inner = func(ctx Context) {}
	}
	return Details(Atr.
		Class(class).
		Id(fmt.Sprint("node-", id)),
		Summary(Atr,
			Span(Atr.
				Add("data-id", fmt.Sprint(id)).
				Class("outline-node-content"), Str(content)),
			nodeInfo(n)),
		inner)
}
