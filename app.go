package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

// This app is going to have a /shared route, everything else is going to be admin stuff
func Route(r chi.Router) {
	// Outline root, show all top-level nodes.
	r.Get("/", HomePage)

	// The HTML page for a node
	r.Get("/node/{id}/page", ShowNodePage)
	r.Get("/node/{id}/page/", ShowNodePage)
	// The data for a node, without the surrounding page bits
	r.Get("/node/{id}/content", ShowNodeContent)
	r.Get("/node/{id}/content/", ShowNodeContent)

	// Create a new node with the information given
	r.Post("/node/create", CreateNodeReq)
	r.Post("/node/create/", CreateNodeReq)

	// Edit node
	r.Post("/node/{id}/edit", UpdateNode)

	// Reparent a node
	r.Post("/node/{id}/reparent", ReparentNode)

	// Remove the node
	r.Delete("/node/{id}/", DeleteNode)

	// Reorder the given list of nodes
	r.Post("/nodes/reorder", ReorderNodes)

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	nodes, err := GetRootNodes()
	die(err)
	HomePageView(w, nodes)
}

// Render a page with a node as it's root
func ShowNodePage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	die(err)
	nodes, err := GetNodesUnder(id)
	die(err)
	HomePageView(w, nodes)
}

// Update a given node's content
// TODO: Also detect and update the meta attributes of the node, if present
func UpdateNode(w http.ResponseWriter, r *http.Request) {
	die(r.ParseForm())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	die(err)
	content := r.PostFormValue("content")
	die(UpdateNodeContent(id, content))
	http.Redirect(w, r, fmt.Sprint("/node/", id, "/page/"), 301)
}

// Only show the HTML fragment (or JSON data, though that's not going to be the starting point)
func ShowNodeContent(w http.ResponseWriter, r *http.Request) {
}

// The fragment for creating a node
func CreateNodeReq(w http.ResponseWriter, r *http.Request) {
	die(r.ParseForm())
	content := r.PostFormValue("content")
	parentIdStr := r.PostFormValue("parent_id")
	outlineOrderStr := r.PostFormValue("outline_order")

	parentId, err := strconv.ParseInt(parentIdStr, 10, 64)
	die(err)
	outlineOrder, err := strconv.ParseInt(outlineOrderStr, 10, 64)
	die(err)
	id, err := CreateNode(OutlineNodeTest{
		ParentId:     parentId,
		Content:      content,
		OutlineOrder: outlineOrder,
		Meta:         "",
	})
	die(err)
	// TODO: Return node later
	// If everything succeeds, write out the id into the response
	fmt.Fprintf(w, "%v", id)
}

func ReparentNode(w http.ResponseWriter, r *http.Request) {
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
}

func ReorderNodes(w http.ResponseWriter, r *http.Request) {
}
