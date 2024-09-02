package parse

import (
	"os"
	"testing"
)

func TestOne(t *testing.T) {
	query := `CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY,name TEXT NOT NULL);`
	fileData, err := os.ReadFile("./table.go")
	if err != nil {
		t.Fatalf("reading table.go failed, %s", err.Error())
	}

	parser := NewParser(fileData)
	sqlStmt := parser.Parse()

	if sqlStmt != query {
		t.Fatalf("TestOne failed, wanted %s got %s", query, sqlStmt)
	}
}
