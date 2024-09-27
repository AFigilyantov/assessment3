package entity

// UserAccount - db schema
type UserAccount struct {
	Username  string `json:"username" validate:"required,min=1,max=20"`
	Password  string `json:"password" validate:"required,min=1"`
	Email     string `json:"email" validate:"required,email"`
	CreatedAt string
}
