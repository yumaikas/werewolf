package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Werewolf is an outliner that works based on SQLite and GoLua")
	TestInitDb()
	TestDbEX()

	r := chi.NewRouter()
	Route(r)
	// TestDb()
}

// This app is going to have a /shared route, everything else is going to be admin stuff
func Route(r chi.Router) {
	// Outline root, show all top-level nodes.
	r.Get("/", HomePage)

	// The HTML page for a node
	r.Get("/node/{id}/page", ShowNodePage)
	// The data for a node, without the surrounding page bits
	r.Get("/node/{id}/content", ShowNodeContent)

	// Create a new node with the information given
	r.Post("/node/create", CreateNodeForm)

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
}

// Only show the HTML fragment (or JSON data, though that's not going to be the starting point)
func ShowNodeContent(w http.ResponseWriter, r *http.Request) {
}

// The fragment for creating a node
func CreateNodeForm(w http.ResponseWriter, r *http.Request) {
}

func ReparentNode(w http.ResponseWriter, r *http.Request) {
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
}

func ReorderNodes(w http.ResponseWriter, r *http.Request) {
}

// Do I even need this structure for this, if I get some good fallback funcitons? I don't think so? I'll see over time.
// Looks like
type OutlineNodeTest struct {
	Id, ParentId int64
	Title        string
	Content      string
	Meta         string
	OutlineOrder int64
	Created      time.Time
	Updated      time.Time
	Deleted      time.Time
}

/*
Going to use an adjacency list for my database.
*/

// Operations I need to support
/*
Re-ordering nodes in a subtree
Inserting nodes
Removing nodes

Executing Lua scripts when nodes are interacted with.
*/
