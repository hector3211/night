package parse

type Users struct {
	ID   int    `orm:"primary_key"`
	Name string `orm:"nullable"`
	// Email         string `orm:"unique"`
	// EmailVerified bool   `orm:"nullable"`
}
