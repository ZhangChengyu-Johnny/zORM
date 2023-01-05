package test

type User struct {
	Name        string `zorm:"PRIMIADY KEY"`
	Age         int
	somePrivate string
}
