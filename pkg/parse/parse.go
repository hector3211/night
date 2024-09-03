package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"night/cmd/flags"
	"reflect"
	"strings"
)

// table example writin in Go
//
//	type Users struct {
//	    ID int `orm:"primary_key"`
//	    Name string
//	    Email string `orm:"unique"`
//	    EmailVerified bool `orm:"nullable"`
//	}
type Table struct {
	Name   string
	Fields [][]string
}

type Parser struct {
	Driver       flags.DataBaseDriver
	fileContents []byte
	Tables       []Table
}

func NewParser(driver flags.DataBaseDriver, contents []byte) *Parser {
	return &Parser{
		Tables:       make([]Table, 0),
		fileContents: contents,
		Driver:       driver,
	}
}

func (p Parser) mapToSql(goType string) string {
	if p.Driver == flags.SQLITE {
		switch goType {
		case "int":
			return "INTEGER"
		case "string":
			return "TEXT"
		case "bool":
			return "BOOL"
		default:
			return "TEXT"
		}
	}
	if p.Driver == flags.POSTGRES {
		// Postgres driver
		switch goType {
		case "int":
			return "INT"
		case "string":
			return "TEXT"
		case "bool":
			return "BOOL"
		default:
			return "TEXT"
		}
	}

	return ""
}

func (p Parser) parseTag(tag string) []string {
	var attributes []string
	tagParts := strings.Split(tag, ",")
	for _, part := range tagParts {
		part = strings.TrimSpace(part)
		switch part {
		case "primary_key":
			attributes = append(attributes, "PRIMARY KEY")
		case "unique":
			attributes = append(attributes, "UNIQUE")
		case "nullable":
			attributes = append(attributes, "NULL")
		case "notnull":
			attributes = append(attributes, "NOT NULL")
		}
	}
	return attributes
}

func (p Parser) generateSql() string {
	var query strings.Builder
	for i := 0; i < len(p.Tables); i++ {
		currTable := p.Tables[i]

		query.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", strings.ToLower(currTable.Name)))
		for idx, fields := range currTable.Fields {
			query.WriteString(fmt.Sprintf("%s", strings.Join(fields, " ")))
			if idx < len(currTable.Fields)-1 {
				query.WriteString(",")
			}
		}
		query.WriteString(")")
		query.WriteString(";")

		if i < len(p.Tables)-1 {
			query.WriteString(" ")
		}

	}
	return query.String()
}

func (p *Parser) Parse() (query string, error error) {
	fset := token.NewFileSet()

	//parse
	node, err := parser.ParseFile(fset, "", string(p.fileContents), parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed reading file %s", err.Error())
	}

	// walk
	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		// We're only interested in structs
		structType, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		table := Table{
			Name: ts.Name.Name,
		}

		for _, field := range structType.Fields.List {
			var fieldInfo []string
			for _, name := range field.Names {
				fieldInfo = append(fieldInfo, strings.ToLower(name.Name))
			}
			fieldInfo = append(fieldInfo, p.mapToSql(fieldType(field.Type)))
			// fmt.Printf("Field: %s, Type: %s", name.Name, fieldType(field.Type))

			if field.Tag != nil {
				tag := reflectStructTag(field.Tag.Value)
				ormTag := tag.Get("orm")
				if ormTag != "" {
					// fmt.Printf(", ORM Tag: %s", ormTag)
					fieldInfo = append(fieldInfo, strings.Join(p.parseTag(ormTag), " "))
				}
			}
			table.Fields = append(table.Fields, fieldInfo)
		}
		p.Tables = append(p.Tables, table)
		return true
	})
	return p.generateSql(), nil
}

// extract the type as a string
func fieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fieldType(t.X) + "." + t.Sel.Name
	default:
		return ""
	}
}

// clean up and parse struct tags
func reflectStructTag(tag string) reflect.StructTag {
	// Remove the backticks from the struct tag
	tag = strings.Trim(tag, "`")
	return reflect.StructTag(tag)
}
