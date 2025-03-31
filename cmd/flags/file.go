package flags

import (
	"fmt"
	"os"
	"path/filepath"
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

func (f *File) Set(value string) error {
	if _, err := os.Stat(value); os.IsNotExist(err) {
		return fmt.Errorf("file does not exists: %s", err)
	}

	ext := filepath.Ext(value)
	if ext == "" {
		return fmt.Errorf("file has no extension, but be on of: %s", strings.Join(AllowedTypes, ","))
	}

	fileType := strings.ToLower(ext[1:])
	for _, allowed := range AllowedTypes {
		if allowed == fileType {
			*f = File(value)
			return nil
		}
	}
	return fmt.Errorf("Provide a valid path")
}
