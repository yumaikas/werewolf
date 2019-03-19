package main

import (
	"fmt"
	"github.com/go-chi/chi"
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

func Route(r chi.Router) {

}

// Do I even need this structure for this, if I get
type OutlineNode struct {
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

Executing Lua scripts when nodes interacted with.
*/
