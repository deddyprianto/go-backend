package models


type User struct{
	ID int `db:"id"`
	Username string `db:"username"`
	Email string `db:"email"`
	Password string `db:"password"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}