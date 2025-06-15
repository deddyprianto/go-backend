package models

import (
	"api-garuda/pkg/helper"
	"database/sql"
	"fmt"
	"time"
)

type User struct{
	ID string `db:"id"`
	Username string `db:"username"`
	Email string `db:"email"`
	Password string `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserLogin struct{
	ID uint `db:"id"`
	Name string `db:"name"`
	Email string `db:"email"`
	Password string `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type UserProfile struct {
    ID        uint      `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}


type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token string `json:"token,omitempty"`
	User *UserLogin `json:"user,omitempty"`
}
// Error ini terjadi karena kita perlu menangani konversi tipe data untuk kolom created_at dan updated_at secara eksplisit ketika melakukan scanning data dari database. MySQL mengembalikan nilai datetime sebagai []uint8, dan kita perlu mengkonversinya ke time.Time.
// Berikut cara memperbaikinya dengan menambahkan method Scan custom untuk struct User:

func (u *User) Scan(rows *sql.Rows) error {
	var createdAt [] uint8
	var updatedAt [] uint8

	err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password,&createdAt, &updatedAt)
	if err != nil{
		return err
	}

	u.CreatedAt, err = helper.Converter(createdAt)

	if err != nil{
		return fmt.Errorf("gagal mengkonversi created_at: %v", err)
	}

	u.UpdatedAt, err = helper.Converter(updatedAt)
	if err != nil{
		return fmt.Errorf("gagal mengkonversi updated_at: %v", err)
	}

	return nil

}
