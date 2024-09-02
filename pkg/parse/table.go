package parse

type Users struct {
	ID   int    `orm:"primary_key"`
	Name string `orm:"notnull"`
}
