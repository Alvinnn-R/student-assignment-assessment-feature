package model

type User struct {
	Model
	Name     string
	Email    string
	Password string
	Role     string
}
