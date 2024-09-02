package flags

import (
	"fmt"
	"strings"
)

type SeedLanguage string

const (
	UNKNWON SeedLanguage = "UNKNWON"
	GOLANG  SeedLanguage = "go"
	SQL     SeedLanguage = "sql"
)

var DefaultSeedLanguage = SQL

// func (s SeedLanguage) String() string {
// 	switch s {
// 	case GOLANG:
// 		return "Go"
// 	case SQL:
// 		return "SQL"
// 	default:
// 		return "SQL"
// 	}
// }

var AllowedFileTypes = []string{"go", "sql"}

func (d SeedLanguage) String() string {
	return string(d)
}

func (f *SeedLanguage) Type() string {
	return "seedLanguage"
}

func (f *SeedLanguage) Set(value string) error {
	for _, fileType := range AllowedFileTypes {
		if fileType == value {
			*f = SeedLanguage(value)
			return nil
		}
	}
	return fmt.Errorf("Allowed seed languages: %s", strings.Join(AllowedFileTypes, ","))
}
