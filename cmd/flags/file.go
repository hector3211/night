package flags

import "fmt"

type File string

func (d File) String() string {
	return string(d)
}

func (f *File) Type() string {
	return "path"
}

func (f *File) Set(value string) error {
	for _, driver := range AllowedDbDrivers {
		if driver == value {
			*f = File(value)
			return nil
		}
	}
	return fmt.Errorf("Provide a valid path")
}
