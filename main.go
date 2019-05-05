package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

var runTests = flag.Bool("test", false, "Set to run tests")
var runServer = flag.Bool("serve", false, "Set to run web server")

func main() {
	flag.Parse()
	fmt.Println("Werewolf is an outliner that works based on SQLite and GoLua")
	// TODO: Create a real db init function
	TestInitDb()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}
	if *runTests {
		// Throw testing behind a CLI flag
		TestDbEX()
	}

	if *runServer {
		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		// TODO: Build a custom 500 page recoverer at some point.
		r.Use(middleware.Recoverer)
		Route(r)

		err := http.ListenAndServe(":4242", r)
		fmt.Println(err)
	}
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
