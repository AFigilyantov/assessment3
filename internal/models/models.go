package models

// TODO transfer to model packege
type UserAccount struct {
	UserName string `json:"username" validate:"required,min=1,max=20"`
	Password string `json:"password" validate:"required,min=1"`
	Email    string `json:"email" validate:"required,email"`
}

type Response struct {
	Error       string
	UserAccount UserAccount
}
