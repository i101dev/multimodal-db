package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --------------------------------------------------------------------
// --------------------------------------------------------------------

var db *gorm.DB

// --------------------------------------------------------------------
// --------------------------------------------------------------------

type User struct {
	gorm.Model
	UUID   string `json:"uuid"`
	Name   string `gorm:"uniqueIndex" json:"name"`
	Skills Skills `gorm:"type:jsonb" json:"skills"`
}
type Skills []Skill
type Skill struct {
	Type  string `json:"type"`
	Level int    `json:"level"`
}

func (s Skills) Value() (driver.Value, error) {
	return json.Marshal(s)
}
func (s *Skills) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, s)
	}
	return errors.New("unsupported data type for scanning into Skills")
}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

func ConnectDB() {

	dbUser := os.Getenv("DB_POSTGRES_USER")
	dbPass := os.Getenv("DB_POSTGRES_PASS")
	dbName := os.Getenv("DB_POSTGRES_NAME")
	dbHost := os.Getenv("DB_POSTGRES_HOST")
	dbPort := os.Getenv("DB_POSTGRES_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" {
		log.Fatal("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	d, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	} else {
		db = d
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Println("Error initializing [models/users.go]")
		log.Fatal(err)
	}
}

func CreateUser(r *http.Request) (*User, error) {

	requestBody := User{
		UUID:   uuid.New().String(),
		Skills: []Skill{},
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	if requestBody.Name == "" {
		return nil, fmt.Errorf("invalid [name]")
	}

	// --------------------------------------------------------------
	// Business logic -----------------------------------------------
	result := db.Create(&requestBody)
	if result.Error != nil {
		return nil, result.Error
	}

	return &requestBody, nil
}

func GetAllUsers(r *http.Request) (*[]User, error) {

	allUsers := []User{}

	result := db.Find(&allUsers)

	if result.Error != nil {
		return &allUsers, result.Error
	}

	if result.RowsAffected == 0 {
		return &allUsers, fmt.Errorf("no users yet")
	}

	return &allUsers, nil
}

func GetUserByID(r *http.Request) (*User, error) {

	user, err := getUserFromUUID(r)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(r *http.Request) (*User, error) {

	requestBody, err := getRequestBody(r)
	if err != nil {
		return nil, fmt.Errorf("error getting request body")
	}

	// ------------------------------------------------------------------------------
	// Fetch user from the database -------------------------------------------------
	userData, err := getUserData(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error getting user data")
	}

	// --------------------------------------------------------------
	// Update user fields -------------------------------------------
	if requestBody.Name != "" {
		userData.Name = requestBody.Name
	}

	// --------------------------------------------------------------
	// Business logic -----------------------------------------------
	if err := db.Save(&userData).Error; err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return userData, nil
}

func DeleteUser(r *http.Request) error {

	user, err := getUserFromUUID(r)

	if err != nil {
		return err
	}

	if err := db.Delete(&user).Error; err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

func getRequestBody(r *http.Request) (*User, error) {

	var requestBody User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	if requestBody.UUID == "" {
		return nil, fmt.Errorf("invalid user UUID")
	}

	return &requestBody, nil
}

func getUserFromUUID(r *http.Request) (*User, error) {

	requestBody, err := getRequestBody(r)
	if err != nil {
		return nil, fmt.Errorf("error getting request body")
	}

	userData, err := getUserData(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error getting user data")
	}

	return userData, nil
}

func getUserData(requestBody *User) (*User, error) {
	user := User{}
	if err := db.Where("uuid = ?", requestBody.UUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &user, nil
}
