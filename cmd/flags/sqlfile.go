package flags

import "fmt"

type SqlFile string

func (d SqlFile) String() string {
	return string(d)
}

func (f *SqlFile) Type() string {
	return "path"
}

func (f *SqlFile) Set(value string) error {
	for _, driver := range AllowedDbDrivers {
		if driver == value {
			*f = SqlFile(value)
			return nil
		}
	}
	return fmt.Errorf("Provide a valid path")
}
