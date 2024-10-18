package model

type UserPtr *User

type User struct {
	baseModel
	Name string
	Age  int
}
