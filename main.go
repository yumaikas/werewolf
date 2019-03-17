package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Werewolf is an outliner that works based on SQLite and GoLua")
	InitDb()
	TestDb()
}

// Do I even need this structure for this, if I get
type OutlineNode struct {
	//
	Id, ParentId int64
	Title        string
	Content      string
	Meta         string
	Created      time.Time
	Updated      time.Time
	Deleted      time.Time
}

// Insert the node at the given IDX in parent, or at the end if the idx is too large
func Insert(parent, child *OutlineNode, idx int) error {
	return nil

}

// Remove the node with the ID given, using idx
// as a hint for where check for the node, reverting to a search if the ID doesn't match
func Remove(parent *OutlineNode, id int64, idx int) error {
	return nil
}

//
//
func Reorder(parent *OutlineNode, new_order_ids []int) error {
	return nil
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
