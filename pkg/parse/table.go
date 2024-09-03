package parse

type Users struct {
	ID   int    `night:"primary_key"`
	Name string `night:"nullable"`
}

type Orders struct {
	OrderID int `night:"primary_key"`
	UserID  int `night:"unique"`
	Amount  string
}
