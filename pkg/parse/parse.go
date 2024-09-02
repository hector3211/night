package parse

import (
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

func (p *Parser) Parse() {

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
}
