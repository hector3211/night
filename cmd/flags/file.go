package flags

import (
	"fmt"
	"strings"
)

type File string

var AllowedTypes = []string{"go", "sql"}

func (d File) String() string {
	return string(d)
}

func (f *File) Type() string {
	return "path"
}

// TODO: fix file path problems
func (f *File) Set(value string) error {
	fileType := strings.Split(value, ".")[1]
	for _, file := range AllowedTypes {
		if file == fileType {
			*f = File(value)
			return nil
		}
	}
	return fmt.Errorf("Provide a valid path")
}
