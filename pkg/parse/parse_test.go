package parse

import (
	"github.com/hector3211/night/cmd/flags"
	"os"
	"testing"
)

func TestOne(t *testing.T) {
	query := `CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY,name TEXT NULL);`
	fileData, err := os.ReadFile("./table.go")
	if err != nil {
		t.Fatalf("reading table.go failed, %s", err.Error())
	}

	parser := NewParser(flags.SQLITE, fileData)
	sqlStmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parsing failed %s", err.Error())
	}

	if sqlStmt != query {
		t.Fatalf("TestOne failed, wanted %s got %s", query, sqlStmt)
	}
}

func TestTwo(t *testing.T) {
	query := `CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY,name TEXT NOT NULL); CREATE TABLE IF NOT EXISTS orders (orderid INT PRIMARY KEY,userid INT UNIQUE,amount TEXT);`
	fileData, err := os.ReadFile("./table.go")
	if err != nil {
		t.Fatalf("reading table.go failed, %s", err.Error())
	}

	parser := NewParser(flags.SQLITE, fileData)
	sqlStmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parsing failed %s", err.Error())
	}

	if sqlStmt != query {
		t.Fatalf("TestTwo failed, wanted %s got %s", query, sqlStmt)
	}
}
