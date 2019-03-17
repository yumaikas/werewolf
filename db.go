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
		(2, "A.A", 1), 
		(2, "A.B", 2), 
		(1, "B", 2),
		(5, "B.A", 1), 
		(5, "B.B", 2);
	`)
	nodes, err := GetNodesUnder(1)

	if err != nil {
		fmt.Println(err)
		return
	}
	for _, n := range nodes {
		fmt.Print(strings.Repeat("*", n.RelativeDepth+1))
		fmt.Println(" " + n.Content.String)
	}
}

func GetNodesUnder(id int64) ([]OutlineNodeDB, error) {
	results := make([]OutlineNodeDB, 0)

	err := db.Select(&results, `
	WITH RECURSIVE nodes(Id, ParentId, depth) as (
		Select Id, ParentId, 0 from Outline where Id = ?
		UNION ALL
		Select Outline.Id, Outline.ParentId, nodes.depth+1
	    from Outline  
	    JOIN nodes ON Outline.ParentId = nodes.Id
	Order By 3 DESC
	) 
	Select 
	    Outline.Id as Id, 
		Outline.ParentId as ParentId, 
		Nodes.Depth as Depth,
		Outline.Content as Content,
		Outline.Meta as Meta,
		Outline.Created as Created,
		Outline.Updated as Updated,
		Outline.Deleted as Deleted
	from Nodes
	INNER JOIN Outline on Nodes.Id = Outline.Id
	`, id)
	return results, err
}

type OutlineNodeDB struct {
	Id       int64         `db:"Id"`
	ParentId sql.NullInt64 `db:"ParentId"`
	// The relative depth of this node from the parent of the current query
	RelativeDepth int `db:"Depth"`

	Content sql.NullString `db:"Content"`
	Meta    sql.NullString `db:"Meta"`
	Created sql.NullInt64  `db:"Created"`
	Updated sql.NullInt64  `db:"Updated"`
	Deleted sql.NullInt64  `db:"Deleted"`
}

// Hrm...
func CreateNode() {
}
