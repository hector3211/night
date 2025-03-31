package flags

import (
	"fmt"
	"strings"
)

type DataBaseDriver string

const (
	SQLITE   DataBaseDriver = "sqlite3"
	POSTGRES DataBaseDriver = "postgres"
)

var AllowedDbDrivers = []string{"sqlite3", "postgres"}

func (f DataBaseDriver) String() string {
	return string(f)
}

func (f *DataBaseDriver) Type() string {
	return "driver"
}

func (f *DataBaseDriver) Set(value string) error {
	for _, driver := range AllowedDbDrivers {
		if driver == value {
			*f = DataBaseDriver(value)
			return nil
		}
	}
	return fmt.Errorf("Allowed Db drivers: %s", strings.Join(AllowedDbDrivers, ","))
}
