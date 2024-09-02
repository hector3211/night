package parse

import (
	"night/cmd/flags"
	"regexp"
	"strings"
)

// table example writin in Go
//
//	type Users struct {
//	    ID night.Int `orm:"primary_key"`
//	    Name night.String
//	    Email night.VarChar `orm:"unique"`
//	    EmailVerified night.Bool `orm:"nullable"`
//	}
type Table struct {
	Name   string
	Fields []map[string]string
}

type Parser struct {
	Driver       flags.DataBaseDriver
	fileContents []byte
	SqlQuery     string
	Tables       []Table
}

func NewParser() *Parser {
	return &Parser{
		Tables: make([]Table, 0),
	}
}

func (p *Parser) SetFileContents(contents []byte) {
	p.fileContents = contents
}

func (p Parser) mapToSql(goType string) string {
	// NOTE: add sqlite and postgres types
	// sqlite types
	// if p.Driver == flags.SQLITE {
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
	// }
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
			// default:
			// Handle other tags if necessary
		}
	}
	return attributes
}

func (p Parser) GenerateSql() string {
	return ""
}

func (p *Parser) Parse() (query string) {

	structReg := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	fieldReg := regexp.MustCompile(`(\w+)\s+(\w+(\.\w+)*)\s*(?:` + "`" + `([^` + "`" + `]*)` + "`" + `)?`)

	structMatches := structReg.FindAllStringSubmatch(string(p.fileContents), -1)
	for _, match := range structMatches {

		structName := match[1]
		fields := match[2]

		var fieldList []map[string]string

		for _, field := range strings.Split(fields, "\n") {
			field = strings.TrimSpace(field)

			if field == "" {
				continue
			}

			fieldMatch := fieldReg.FindStringSubmatch(field)
			if fieldMatch != nil {
				fieldInfo := make(map[string]string, 0)
				fieldInfo["name"] = fieldMatch[1]
				fieldInfo["type"] = fieldMatch[2]

				// TODO: Look over this and figure out types
				tag := fieldMatch[4]
				if tag != "" {
					nightTag := strings.Split(tag, " ")
					for _, part := range nightTag {
						if strings.HasPrefix(part, "orm:") {
							fieldInfo["orm"] = strings.TrimPrefix(part, "orm:")
						}
					}
				}
				fieldList = append(fieldList, fieldInfo)
			}
			p.Tables = append(p.Tables, Table{Name: structName, Fields: fieldList})
		}
	}
	return ""
}
