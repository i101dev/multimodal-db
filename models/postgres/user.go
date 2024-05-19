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

	"github.com/i101dev/multimodal-db/util"
)

// --------------------------------------------------------------------
// --------------------------------------------------------------------

var db *gorm.DB

// --------------------------------------------------------------------
// --------------------------------------------------------------------

type User struct {
	gorm.Model
	UUID     string `json:"uuid"`
	Name     string `gorm:"uniqueIndex" json:"name"`
	Location string `json:"location"`
	Skills   Skills `gorm:"type:jsonb" json:"skills"`
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

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	d, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal("\n*** >>> Postgres connection failed:", err)
		return
	}

	db = d

	// ----------------------------------------------------
	// Migrations -----------------------------------------
	//
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal("Error initializing [models/users.go]:", err)
	}
}

func CreateUser(r *http.Request) (*User, error) {

	requestBody, userData, _ := userData_byName(r)

	if userData != nil {
		return nil, fmt.Errorf("name already in play")
	}

	// ----------------------------------------------
	//
	if requestBody.Name == "" {
		return nil, fmt.Errorf("invalid [name]")
	}
	if requestBody.Location == "" {
		return nil, fmt.Errorf("invalid [location]")
	}
	//
	// ----------------------------------------------

	requestBody.UUID = uuid.New().String()
	requestBody.Skills = []Skill{}

	if result := db.Create(&requestBody); result.Error != nil {
		return nil, result.Error
	}

	return requestBody, nil
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

func FindUserByID(r *http.Request) (*User, error) {

	_, userData, err := userData_byUUID(r)

	if err != nil {
		return nil, err
	}

	return userData, nil
}

func UpdateUser(r *http.Request) (*User, error) {

	requestBody, userData, err := userData_byUUID(r)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------
	//
	if requestBody.Name != "" {
		userData.Name = requestBody.Name
	}
	if requestBody.Location != "" {
		userData.Location = requestBody.Location
	}
	//
	// ----------------------------------------------

	if err := db.Save(&userData).Error; err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return userData, nil
}

func DeleteUser(r *http.Request) error {

	_, userData, err := userData_byUUID(r)

	if err != nil {
		return err
	}

	if err := db.Delete(&userData).Error; err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

func userData_byUUID(r *http.Request) (*User, *User, error) {

	var requestBody User

	if err := util.ParseBody(r, &requestBody); err != nil {
		return nil, nil, err
	}

	if requestBody.UUID == "" {
		return nil, nil, fmt.Errorf("invalid user [UUID]")
	}

	userData := &User{}

	if err := db.Where("uuid = ?", requestBody.UUID).First(userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &requestBody, nil, fmt.Errorf("user not found")
		}
		return nil, nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &requestBody, userData, nil
}

func userData_byName(r *http.Request) (*User, *User, error) {

	var requestBody User

	if err := util.ParseBody(r, &requestBody); err != nil {
		return nil, nil, err
	}

	if requestBody.Name == "" {
		return nil, nil, fmt.Errorf("invalid user [Name]")
	}

	userData := &User{}

	if err := db.Where("name = ?", requestBody.Name).First(userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &requestBody, nil, fmt.Errorf("user not found")
		}
		return nil, nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &requestBody, userData, nil
}
