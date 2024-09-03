package parse

type Users struct {
	ID   int    `orm:"primary_key"`
	Name string `orm:"nullable"`
}

type Orders struct {
	OrderID int `orm:"primary_key"`
	UserID  int `orm:"unique"`
	Amount  string
}
