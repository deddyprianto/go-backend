package database

import (
	"api-garuda/pkg/helper"
	"api-garuda/pkg/models"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Response struct{
	Data []models.User `json:"data"`
	Message string     `json:"message"`
}

type ResponseSaveSingle struct{
    Message string `json:"message"`
    Data int `json:"data"`
}

func GetAllUSers(db *sql.DB) (Response, error) {
    query := "SELECT * FROM users"
    rows, err := db.Query(query)
    if err != nil {
        return Response{Message: err.Error()}, fmt.Errorf("gagal eksekusi query: %v", err)
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        err := user.Scan(rows)
        if err != nil {
            return Response{Message: err.Error()}, fmt.Errorf("gagal parsing data: %s", err)
        }
        users = append(users, user)
    }

    if len(users) == 0 {
        return Response{
			Data: []models.User{},
			Message: "data not found",
		},nil
    }

    return Response{
		Data: users,
		Message: "success",
	}, nil
}

func GetUserById(db *sql.DB, id string) (models.User, error) {
    query := "SELECT id, username, email, password, created_at, updated_at , date_modification FROM users WHERE id = ?"
    row := db.QueryRow(query, id)
    
    // Gunakan variabel sementara untuk menyimpan hasil query
    var (
        userId      string
        username    string
        email       string
        password    string
        createdAt  []uint8
        updatedAt  []uint8
        date_modification string
    )
    
    // Lakukan scanning ke variabel sementara
    err := row.Scan(&userId, &username, &email, &password, &createdAt, &updatedAt, &date_modification)
    if err != nil {
        return models.User{}, fmt.Errorf("failed to scan user data: %v", err)
    }
    
    // Konversi byte array ke time.Time
    createdAtTime, err := helper.Converter(createdAt)
    if err != nil {
        return models.User{}, fmt.Errorf("failed to convert created_at time: %v", err)
    }
    
    updatedAtTime, err := helper.Converter(updatedAt)
    if err != nil {
        return models.User{}, fmt.Errorf("failed to convert updated_at time: %v", err)
    }
    
    // Buat dan return struct User yang sudah dikonversi
    return models.User{
        ID:        userId,
        Username: username,
        Email:     email,
        Password: password,
        CreatedAt: createdAtTime,
        UpdatedAt: updatedAtTime,
        DateModification: date_modification,
    }, nil
}


func UpdateUser(db *sql.DB, user models.User, id string) (string, error){
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := "UPDATE users SET username=?, email=?, password=?, created_at=?, updated_at=? WHERE id=?"

	stmt, err := db.Prepare(query)
	if err != nil{
        return "", fmt.Errorf("failed to prepare update query in UpdateUser function: %v", err)
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, id)
	if err != nil{
		return "", fmt.Errorf("failed to update user with ID %v: %v", id, err)
	}

	if affected , err := result.RowsAffected(); err != nil || affected == 0{
        return "", fmt.Errorf("no changes detected when updating user with ID %v", id)
	}
	return id , nil
}

func CreateUser(db *sql.DB, user models.User) (ResponseSaveSingle ,error){
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	query := "INSERT INTO users (username, email, password, created_at, updated_at, date_modification) VALUES (?, ?, ?, ?, ?, ?)"

	stmt , err := db.Prepare(query)
	if err != nil{
		return ResponseSaveSingle{Message: err.Error()}, fmt.Errorf("failed prepare data: %v", err)
	}
	defer stmt.Close()

	result ,err := stmt.Exec(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.DateModification)

	if err != nil{
		return ResponseSaveSingle{Message: err.Error()} , fmt.Errorf("failed insert data: %v", err)
	}

    id, err := result.LastInsertId()
    if err != nil{
        return ResponseSaveSingle{
            Message: err.Error(),
            Data: 0,
        }, err
    }
    return ResponseSaveSingle{
        Message: "success save data with id: " + strconv.Itoa(int(id)),
        Data: int(id),
    }, nil
}

func DeleteUser(db *sql.DB, id string) (int64, error){
	query := "DELETE FROM users WHERE id = ?"
	stmt , err := db.Prepare(query)
	if err != nil{
        return 0, fmt.Errorf("gagal prepare delete query: %v", err)
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(id)
	if err != nil {
        return 0, fmt.Errorf("gagal delete data: %v", err)
	}
	return result.RowsAffected()
}

