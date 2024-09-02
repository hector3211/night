package parse

type Users struct {
	ID            int `orm:"primary_key"`
	Name          string
	Email         string `orm:"unique"`
	EmailVerified bool   `orm:"nullable"`
}
