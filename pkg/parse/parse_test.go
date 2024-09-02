package parse

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestOne(t *testing.T) {
	// query := []byte("type Users struct {\nID int\nName string\nEmail string\nEmailVerified bool\n}")
	fileData, err := os.ReadFile("./table.go")
	if err != nil {
		t.Fatalf("reading table.go failed, %s", err.Error())
	}
	structReg := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	fieldReg := regexp.MustCompile(`(\w+)\s+(\w+(\.\w+)*)\s*(?:` + "`" + `([^` + "`" + `]*)` + "`" + `)?`)

	structMatches := structReg.FindAllStringSubmatch(string(fileData), -1)
	var tables []Table
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
		}
		tables = append(tables, Table{Name: structName, Fields: fieldList})
	}

	if tables[0].Name != "Users" {
		t.Fatalf("failed structName wanted Users got %s", tables[0].Name)
	}

	// expectedFields := []string{"ID", "Name", "Email", "EmailVerified"}
	//
	// for idx, col := range expectedFields {
	// 	if fieldList[idx]["name"] != col {
	// 		t.Fatalf("faileds dont match wanted %s got %s", col, fieldList[idx]["name"])
	// 	}
	// }
}
