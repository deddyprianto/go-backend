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

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Data    []models.User `json:"data"`
	Message string        `json:"message"`
}

type ResponseSaveSingle struct {
	Message string      `json:"message"`
	Data    models.User `json:"data"`
}
type ResponseSaveSingleEmployee struct {
	Message string          `json:"message"`
	Data    models.Employee `json:"data"`
}

type ResponseOnUpdate struct {
	Message string      `json:"message"`
	Data    models.User `json:"data"`
}

// Tambahkan struct Claims untuk menangani token JWT
type Claims struct {
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ExpiresAt      jwt.NumericDate
	StandardClaims jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if c.ExpiresAt.IsZero() {
		return nil
	}
	// Jika current time sudah melebihi expiresAt, token sudah expired
	if time.Now().Unix() > c.ExpiresAt.Unix() {
		return fmt.Errorf("token expired")
	}
	return nil
}

func ValidateToken(token string, secretKey string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token invalid or expired: %v", err)
	}
	return claims, nil
}

func GenerateNewToken(userID uint, email string, secretKey string) (*TokenPair, error) {
	// Access token expiration (15 menit)
	accessTokenExp := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	// Refresh token expiration (1 minggu)
	refreshTokenExp := jwt.NewNumericDate(time.Now().Add(168 * time.Hour))

	// Generate access token
	accessClaims := &Claims{
		UserID: userID,
		Name:   "User",
		Email:  email,
		StandardClaims: jwt.RegisteredClaims{
			ExpiresAt: accessTokenExp,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("gagal generate access token: %v", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ExpiresAt: *refreshTokenExp,
	}).SignedString([]byte(secretKey))

	if err != nil {
		return nil, fmt.Errorf("gagal generate refresh token: %v", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	query := "SELECT id, username, email,date_modification, created_at FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	// Gunakan variabel sementara untuk menyimpan hasil query
	var (
		userId            string
		username          string
		email             string
		date_modification string
		createdAt         []uint8
	)

	// Lakukan scanning ke variabel sementara
	err := row.Scan(&userId, &username, &email, &date_modification, &createdAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to scan user data: %v", err)
	}

	// Konversi byte array ke time.Time
	createdAtTime, err := helper.Converter(createdAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to convert created_at time: %v", err)
	}

	// Buat dan return struct User yang sudah dikonversi
	return models.User{
		ID:        userId,
		Username:  username,
		Email:     email,
		CreatedAt: createdAtTime,
	}, nil
}

func UpdateUser(db *sql.DB, user models.User, id string) (ResponseOnUpdate, error) {
	user.CreatedAt = time.Now()

	query := "UPDATE users SET username=?, email=?, created_at=? WHERE id=?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseOnUpdate{Message: "failed to prepare update query"}, fmt.Errorf("failed to prepare update query in UpdateUser function: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.CreatedAt, id)
	if err != nil {
		return ResponseOnUpdate{Message: fmt.Sprintf("failed to update user with ID %v", id)}, fmt.Errorf("failed to update user with ID %v: %v", id, err)
	}

	if affected, err := result.RowsAffected(); err != nil || affected == 0 {
		return ResponseOnUpdate{Message: fmt.Sprintf("no changes detected when updating user with ID %v", id)}, fmt.Errorf("no changes detected when updating user with ID %v", id)
	}
	return ResponseOnUpdate{Message: "success update user", Data: user}, nil
}

func CreateEmployee(db *sql.DB, employee models.Employee) (ResponseSaveSingleEmployee, error) {
	employee.CreatedAt = time.Now()

	query := "INSERT INTO employees (name, position, created_at, profile_picture) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseSaveSingleEmployee{Message: err.Error()}, fmt.Errorf("failed to prepare insert query: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(employee.Name, employee.Position, employee.CreatedAt, employee.ProfilePicture)
	if err != nil {
		return ResponseSaveSingleEmployee{Message: err.Error()}, fmt.Errorf("failed to insert employee data: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return ResponseSaveSingleEmployee{
			Message: err.Error(),
			Data:    models.Employee{},
		}, err
	}
	return ResponseSaveSingleEmployee{
		Message: "success save data with id: " + strconv.Itoa(int(id)),
		Data:    employee,
	}, nil
}

func CreateUser(db *sql.DB, user models.User) (ResponseSaveSingle, error) {
	user.CreatedAt = time.Now()
	query := "INSERT INTO users (username, email, date_modification, created_at) VALUES (?, ?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseSaveSingle{Message: err.Error()}, fmt.Errorf("failed prepare data: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.DateModification, user.CreatedAt)

	if err != nil {
		return ResponseSaveSingle{Message: err.Error()}, fmt.Errorf("failed insert data: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return ResponseSaveSingle{
			Message: err.Error(),
			Data:    models.User{},
		}, err
	}
	return ResponseSaveSingle{
		Message: "success save data with id: " + strconv.Itoa(int(id)),
		Data: models.User{
			ID:        strconv.Itoa(int(id)),
			Username:  user.Username,
			Email:     user.Email,
			DateModification: user.DateModification,
			CreatedAt: user.CreatedAt,
		},
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
