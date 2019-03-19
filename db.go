package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
		Created int, -- unix timestamp
		Updated int, -- unix timestamp
		Deleted int -- Unix timestamp
	);

	Create Table If Not Exists Scripts (
		Id INTEGER PRIMARY KEY,
		Name text,
		Code text
	);
	`)
}

func TestDb() {
	// Testing only
	db.MustExec(`DROP TABLE IF EXISTS Outline`)
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

func PrintNodesUnder(id int64) {
	nodes, err := GetNodesUnder(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, n := range nodes {
		fmt.Print(strings.Repeat("*", n.RelativeDepth+1))
		fmt.Println(" "+n.Content.String, " ", n.Id)
	}
}

func die(e error) {
	if e != nil {
		panic(e)
	}
}

func TestDbEX() {
	TestDb()
	id, err := CreateNode(OutlineNode{
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
	// TODO: Test remove and reorder later
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
func CreateNode(node OutlineNode) (int64, error) {

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
