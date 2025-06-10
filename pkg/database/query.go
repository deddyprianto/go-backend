package database

import (
	"api-garuda/pkg/models"
	"api-garuda/pkg/helper"
	"database/sql"
	"fmt"
	"time"
)

func GetAllUSers(db *sql.DB) ([]models.User, error){
	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil{
		return nil,fmt.Errorf("gagal eksekusi query: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next(){
		var user models.User
		err := user.Scan(rows)
		if err != nil{
			return nil , fmt.Errorf("gagal parsing data: %s", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserById(db *sql.DB, id string) (models.User, error) {
    query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?"
    row := db.QueryRow(query, id)
    
    // Gunakan variabel sementara untuk menyimpan hasil query
    var (
        userId      string
        username    string
        email       string
        password    string
        createdAt  []uint8
        updatedAt  []uint8
    )
    
    // Lakukan scanning ke variabel sementara
    err := row.Scan(&userId, &username, &email, &password, &createdAt, &updatedAt)
    if err != nil {
        return models.User{}, fmt.Errorf("gagal mengambil data dengan id: %v", err)
    }
    
    // Konversi byte array ke time.Time
    createdAtTime, err := helper.Converter(createdAt)
    if err != nil {
        return models.User{}, fmt.Errorf("gagal mengkonversi created_at: %v", err)
    }
    
    updatedAtTime, err := helper.Converter(updatedAt)
    if err != nil {
        return models.User{}, fmt.Errorf("gagal mengkonversi updated_at: %v", err)
    }
    
    // Buat dan return struct User yang sudah dikonversi
    return models.User{
        ID:        userId,
        Username: username,
        Email:     email,
        Password: password,
        CreatedAt: createdAtTime,
        UpdatedAt: updatedAtTime,
    }, nil
}


func UpdateUser(db *sql.DB, user models.User) (int64,error){
	query := "UPDATE users SET username = ?, email = ?, password = ?, created_at = ?, updated_at = ? WHERE id = ?"
	stmt, err  := db.Prepare(query)

	if err != nil{
		return 0, fmt.Errorf("anda gagal update data : %v", err)
	}
	defer stmt.Close()

	result , err := stmt.Exec(user.Username,user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.ID)

	if err != nil{
		return 0, fmt.Errorf("kau gagal anjing: %v", err)
	}

	return result.RowsAffected()
}


func CreateUser(db *sql.DB, user models.User) (int64 ,error){
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	query := "INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"

	stmt , err := db.Prepare(query)
	if err != nil{
		return 0 , fmt.Errorf("anda gagal save data: %v", err)
	}
	defer stmt.Close()

	result ,err := stmt.Exec(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err != nil{
		return 0 , fmt.Errorf("GAGAL EXECUTE QUERY: %v", err)
	}
	return result.LastInsertId()
}

func DeleteUser(db *sql.DB, id string) (int64, error){
	query := "DELETE FROM users WHERE id = ?"
	stmt , err := db.Prepare(query)
	if err != nil{
		return 0, fmt.Errorf("gagal delete data: %v", err)
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("kau gagal anjing delete data : %v" , err)
	}
	return result.RowsAffected()
}

