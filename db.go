package main

import (
	"database/sql"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

var db *sqlx.DB

func InitDb() {
	var err error
	db, err = sqlx.Open("sqlite3", "outline.db")
	if err != nil {
		panic(err)
	}
}

func TestInitDb() {
	var err error
	db, err = sqlx.Open("sqlite3", "test-outline.db")
	if err != nil {
		panic(err)
	}
}

func CreateDb() {
	db.MustExec(`
	Create Table If Not Exists Outline (
		Id INTEGER PRIMARY KEY,
		ParentId int,
		OutlineOrder int, 
		Content text,
		Meta text,
		IsExpanded int,
		Created int, -- Unix timestamp
		Updated int, -- Unix timestamp
		Deleted int -- Unix timestamp
	);

	Create Table If Not Exists Scripts (
		Id INTEGER PRIMARY KEY,
		Name text,
		Code text,
		Meta text,
		Created int, -- Unix timestamp
		Updated int, -- Unix timestamp
		Deleted int -- Unix timestamp
	);
	`)
}

func TestDb() {
	// Testing only
	db.MustExec(`DROP TABLE IF EXISTS Outline;`)
	CreateDb()
	// Testing Purposes only
	db.MustExec(`
	Insert Into Outline(ParentId, Content, OutlineOrder) values 
		(NULL, "TOP", 0), 
		(1, "A", 1), 
		(2, "A.A", 4), 
		(2, "A.B", 1), 
		(1, "B", 2),
		(5, "B.A", 1), 
		(5, "B.B", 4);
	`)
	PrintNodesUnder(1)
}

/*
What lua bits do I need?

# Events
- I need something to add content to an existing Outline Node when it's loaded?

# REPL
- Gather nodes

# Nodes
- Create custom node types
- Search nodes
- get the descendants of nodes
- Add virtual subnodes
- Add global gadgets

"type:gadget type:custom-node type:global-gadget"
"enabled:false enabled:true"
""

*/

//
func TestScriptDb() {
	CreateScript("", "type:gadget type:cust", `

	`)
}

func PrintNodes(nodes []OutlineNodeDB) {
	for _, n := range nodes {
		fmt.Print(strings.Repeat("*", n.RelativeDepth+1))
		fmt.Println(" "+n.Content.String, " ", n.Id)
	}
}

func PrintNodesUnder(id int64) {
	nodes, err := GetNodesUnder(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintNodes(nodes)
	HomePageView(os.Stdout, nodes)
}

func die(e error) {
	if e != nil {
		panic(e)
	}
}

func TestDbEX() {
	TestDb()
	id, err := CreateNode(OutlineNodeTest{
		ParentId:     2,
		Content:      "A.B.C",
		OutlineOrder: 6,
		Meta:         "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("New Node Id", id)
	die(Reparent(8, 5, 10))
	PrintNodesUnder(1)
	roots, err := GetRootNodes()
	die(err)
	PrintNodes(roots)
	// TODO: Test remove and reorder later
}

func GetRootNodes() ([]OutlineNodeDB, error) {
	results := make([]OutlineNodeDB, 0)

	err := db.Select(&results, `
		Select 
			Id, 
			ParentId, 
			0 as Depth, 
			Outline.OutlineOrder,
			Outline.Content as Content,
			Outline.Meta as Meta,
			Outline.Created as Created,
			Outline.Updated as Updated,
			Outline.Deleted as Deleted
		From Outline where ParentId is null
	`)
	return results, err
}

func GetNodesUnder(id int64) ([]OutlineNodeDB, error) {
	results := make([]OutlineNodeDB, 0)

	err := db.Select(&results, `
	WITH RECURSIVE nodes(
		Id, 
		ParentId, 
		Depth, 
		OutlineOrder,
		Content,
		Meta,
		Created,
		Updated,
		Deleted
	) as (
		Select 
			Id, 
			ParentId, 
			0, 
			Outline.OutlineOrder,
			Outline.Content as Content,
			Outline.Meta as Meta,
			Outline.Created as Created,
			Outline.Updated as Updated,
			Outline.Deleted as Deleted
		From Outline where Id = ?
		UNION ALL
		Select 
			Outline.Id as Id, 
			Outline.ParentId as ParentId, 
			Nodes.Depth+1 as Depth,
			Outline.OutlineOrder as OutlineOrder,
			Outline.Content as Content,
			Outline.Meta as Meta,
			Outline.Created as Created,
			Outline.Updated as Updated,
			Outline.Deleted as Deleted
	    from Outline  
	    JOIN nodes ON Outline.ParentId = nodes.Id 
		Order By 3 DESC, 4 ASC
	) 
	Select * from Nodes
	where Deleted is Null

	`, id)
	return results, err
}

func ScriptContentForName(name string) (string, error) {
	var content string
	err := db.Get(&content, `Select Content from Scripts where name = ?;`, name)
	return content, err
}

type OutlineNodeDB struct {
	Id       int64         `db:"Id"`
	ParentId sql.NullInt64 `db:"ParentId"`
	// The relative depth of this node from the parent of the current query
	RelativeDepth int           `db:"Depth"`
	OutlineOrder  sql.NullInt64 `db:"OutlineOrder"`

	Content sql.NullString `db:"Content"`
	Meta    sql.NullString `db:"Meta"`
	Created sql.NullInt64  `db:"Created"`
	Updated sql.NullInt64  `db:"Updated"`
	Deleted sql.NullInt64  `db:"Deleted"`
}

type OutlineTree struct {
	Self     OutlineNodeDB
	Children []*OutlineTree
}

func (me *OutlineTree) AddChild(c *OutlineTree) {
	if len(me.Children) == 0 {
		me.Children = []*OutlineTree{c}
		return
	}
	me.Children = append(me.Children, c)
}

func inSet(set map[int64]*OutlineNodeDB, id int64) bool {
	_, found := set[id]
	return found
}

// Depends on parent nodes being listed before their children...
// Given a list of nodes, return a list of roots
func NodesToTree(nodes []OutlineNodeDB) []*OutlineTree {
	// First, build a set of all the ids
	ids := make(map[int64]*OutlineNodeDB)
	idList := make([]int64, 0)
	for _, n := range nodes {
		ids[n.Id] = &n
		idList = append(idList, n.Id)
	}
	topLevels := make([]*OutlineTree, 0)
	treeMap := make(map[int64]*OutlineTree)
	for _, n := range nodes {

		branch := OutlineTree{Self: n}
		treeMap[n.Id] = &branch
		// Handle top-level nodes
		if !n.ParentId.Valid || !inSet(ids, n.ParentId.Int64) {
			topLevels = append(topLevels, &branch)
			continue
		}
		// If not a top-level node, add it to it's parent node
		// So, I'm getting an error where this treemap
		// isn't holding one of the parents like it should yet?
		// Perhaps I need an check here, and add the parent node if it doesn't exist?
		// TODO: Try to fix this later, with more food in my body.
		treeMap[n.ParentId.Int64].AddChild(&branch)
	}
	return topLevels
}

func Int64Or(base sql.NullInt64, fallback int64) int64 {
	if base.Valid {
		return base.Int64
	}
	return fallback
}

func StringOr(base sql.NullString, fallback string) string {
	if base.Valid {
		return base.String
	}
	return fallback
}

/*
	Create Table If Not Exists Outline (
		Id INTEGER PRIMARY KEY,
		ParentId int,
		OutlineOrder int,
		Content text,
		Meta text,
		Created int, -- unix timestamp
		Updated int, -- unix timestamp
		Deleted int -- Unix timestamp
	);
*/
// Create a node
func CreateNode(node OutlineNodeTest) (int64, error) {

	res, err := db.Exec(`Insert into Outline(ParentId, OutlineOrder, Content, Meta, Created)
	values (?, ?, ?, ?, strftime('%s', 'now'));`,
		node.ParentId, node.OutlineOrder, node.Content, node.Meta)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}

	return id, nil
}

// Update a node's content and/or meta
// TODO-@meta: Add in meta support here
func UpdateNodeContent(id int64, newContent string) error {
	_, err := db.Exec(`
	Update Outline 
		Set Content = ?, 
		Updated = strftime('now', '%s')
	where Id = ?`, newContent, id)
	return err
}

// Move the current node under a new parent with a given order number.
func Reparent(Id, NewParentId, OrderNum int64) error {
	_, err := db.Exec(`
	Update Outline 
	Set 
		ParentId = ?,
		OutlineOrder = ?
	where Id = ?;
	`, NewParentId, OrderNum, Id)
	if err != nil {
		return err
	}
	return nil
}

// Remove the node with the ID given, using idx
// as a hint for where check for the node, reverting to a search if the ID doesn't match
func Remove(Id int64) error {
	_, err := db.Exec(`
	Update Outline Set
	    Deleted = strftime('%s', 'now')
	from Outline where Id = ?
	`, Id)
	if err != nil {
		return err
	}
	return nil
}

// Re-arrange the nodes to match the new order
func Reorder(parentId int, newOrderIds []int) error {
	for idx, id := range newOrderIds {
		_, err := db.Exec(`
		Update Outline 
			Set OutlineOrder = ? 
		from Outline where Id = ? and ParentId = ?`, idx, id, parentId)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateScript(name, meta, content string) (int64, error) {
	results, _ := db.Exec(`
	Insert Into Scripts(Name, Code, Meta, Created) 
	Values (?, ?, ?, strftime('%s', 'now');
	`, name, content, meta)
	return results.LastInsertId()
}

func UpdateScript(id int64, name, content string) (int64, error) {
	results, _ := db.Exec(`
	Update Scripts
	   Set 
		   Name = ?,
	       Code = ?,
		   Updated = strftime('%s', 'now')
	   where Id = ?;
	`, name, content, id)
	return results.LastInsertId()
}
