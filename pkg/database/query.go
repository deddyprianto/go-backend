package database

import (
	"api-garuda/internal/middleware"
	"api-garuda/pkg/helper"
	"api-garuda/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Data    []models.User `json:"data"`
	Message string        `json:"message"`
}

type ResponseSaveSingle struct {
	Message string `json:"message"`
	Data    int    `json:"data"`
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
			Data:    []models.User{},
			Message: "data not found",
		}, nil
	}

	return Response{
		Data:    users,
		Message: "success",
	}, nil
}

func GetUserById(db *sql.DB, id string) (models.User, error) {
	query := "SELECT id, username, email, password, created_at, updated_at , date_modification FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	// Gunakan variabel sementara untuk menyimpan hasil query
	var (
		userId            string
		username          string
		email             string
		password          string
		createdAt         []uint8
		updatedAt         []uint8
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
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: createdAtTime,
		UpdatedAt: updatedAtTime,
	}, nil
}

func UpdateUser(db *sql.DB, user models.User, id string) (string, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := "UPDATE users SET username=?, email=?, password=?, created_at=?, updated_at=? WHERE id=?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare update query in UpdateUser function: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, id)
	if err != nil {
		return "", fmt.Errorf("failed to update user with ID %v: %v", id, err)
	}

	if affected, err := result.RowsAffected(); err != nil || affected == 0 {
		return "", fmt.Errorf("no changes detected when updating user with ID %v", id)
	}
	return id, nil
}

func CreateUser(db *sql.DB, user models.User) (ResponseSaveSingle, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := "INSERT INTO users (username, email, password, created_at, updated_at, date_modification) VALUES (?, ?, ?, ?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseSaveSingle{Message: err.Error()}, fmt.Errorf("failed prepare data: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return ResponseSaveSingle{Message: err.Error()}, fmt.Errorf("failed insert data: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return ResponseSaveSingle{
			Message: err.Error(),
			Data:    0,
		}, err
	}
	return ResponseSaveSingle{
		Message: "success save data with id: " + strconv.Itoa(int(id)),
		Data:    int(id),
	}, nil
}

func DeleteUser(db *sql.DB, id string) (int64, error) {
	query := "DELETE FROM users WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("gagal prepare delete query: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("gagal delete data: %v", err)
	}
	return result.RowsAffected()
}

func Register(db *sql.DB, req models.RegisterRequest) (*models.AuthResponse, error) {
	// check if user already exist
	var existingUser models.UserLogin
	query := "SELECT id FROM userslogin WHERE email = ?"
	row := db.QueryRow(query, req.Email)
	err := row.Scan(&existingUser.ID)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// insert new user
	result, err := db.Exec("INSERT INTO userslogin (name,email,password, created_at, updated_at) VALUES (?,?,?,NOW(),NOW())", req.Name, req.Email, string(hashPassword))

	if err != nil {
		return nil, errors.New("failed to create user")
	}

	userID, _ := result.LastInsertId()

	// generate JWT Token
	tokenPair, err := middleware.GenerateTokenPair(uint(userID), req.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// return response
	user := &models.UserLogin{
		ID:    uint(userID),
		Name:  req.Name,
		Email: req.Email,
	}

	return &models.AuthResponse{
		Message: "user Register success",
		Token:   tokenPair,
		User:    user,
	}, nil
}

func Login(db *sql.DB, req models.LoginRequest) (*models.AuthResponse, error) {
	// dapatkan user dari database
	var user models.UserLogin
	query := "SELECT id, name, email, password FROM userslogin WHERE email = ?"
	row := db.QueryRow(query, req.Email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	// return response don't include password
	user.Password = ""
	return &models.AuthResponse{
		Message: "success login",
		User:    &user,
	}, nil
}

func parseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",           // MySQL DATETIME
		"2006-01-02T15:04:05Z",          // ISO 8601 UTC
		"2006-01-02T15:04:05.000Z",      // ISO 8601 dengan milliseconds
		"2006-01-02T15:04:05-07:00",     // ISO 8601 dengan timezone
		"2006-01-02T15:04:05.000-07:00", // ISO 8601 dengan milliseconds dan timezone
		time.RFC3339,                    // RFC3339
		time.RFC3339Nano,
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateStr)
}

func GetProfile(db *sql.DB, userID uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	query := "SELECT id,name,email,created_at,updated_at FROM userslogin WHERE id= ?"
	row := db.QueryRow(query, userID)

	var (
		id        uint
		name      string
		email     string
		createdAt []uint8
		updatedAt []uint8
	)

	// Gunakan sql.NullTime untuk menangani created_at dan updated_at
	err := row.Scan(&id, &name, &email, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	createdAtStr := string(createdAt)
	parsedCreateAt, err := parseDateTime(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %v", err)
	}

	updatedAtStr := string(updatedAt)
	parsedUpdatedAt, err := parseDateTime(updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %v", err)
	}

	profile.ID = id
	profile.Name = name
	profile.Email = email
	profile.CreatedAt = parsedCreateAt
	profile.UpdatedAt = parsedUpdatedAt

	return &profile, nil
}
