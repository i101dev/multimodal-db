package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/google/uuid"
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
	Name     string `gorm:"type:varchar(255);uniqueIndex" json:"name"`
	Location string `json:"location"`
	Skills   Skills `gorm:"type:json" json:"skills"`
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

	dbUser := os.Getenv("DB_MYSQL_USER")
	dbPass := os.Getenv("DB_MYSQL_PASSWORD")
	dbName := os.Getenv("DB_MYSQL_DATABASE")
	dbHost := os.Getenv("DB_MYSQL_HOST")
	dbPort := os.Getenv("DB_MYSQL_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("incomplete database connection parameters")
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	d, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal("\n*** >>> MySQL connection failed:", err)
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
		return nil, fmt.Errorf("name already in use")
	}

	// ----------------------------------------------
	//
	if requestBody.Name == "" {
		return nil, fmt.Errorf("invalid [name]")
	}
	if requestBody.Location == "" {
		return nil, fmt.Errorf("invalid [location]")
	}
	// ----------------------------------------------

	requestBody.UUID = uuid.New().String()
	requestBody.Skills = []Skill{}

	if result := db.Create(requestBody); result.Error != nil {
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
	} else if requestBody.Name == "" && requestBody.Location == "" {
		return nil, fmt.Errorf("nothing to update")
	}

	// ----------------------------------------------
	//
	if requestBody.Name != "" {
		userData.Name = requestBody.Name
	}
	if requestBody.Location != "" {
		userData.Location = requestBody.Location
	}
	// ----------------------------------------------

	if err := db.Save(userData).Error; err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return userData, nil
}

func DeleteUser(r *http.Request) error {

	_, userData, err := userData_byUUID(r)

	if err != nil {
		return err
	}

	if err := db.Delete(userData).Error; err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

func userData_byUUID(r *http.Request) (*User, *User, error) {

	var reqBody User

	if err := util.ParseBody(r, &reqBody); err != nil {
		return &reqBody, nil, err
	}

	if reqBody.UUID == "" {
		return &reqBody, nil, fmt.Errorf("invalid user [UUID]")
	}

	userData := &User{}

	if err := db.Where("uuid = ?", reqBody.UUID).First(userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &reqBody, nil, fmt.Errorf("user not found")
		}
		return &reqBody, nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &reqBody, userData, nil
}

func userData_byName(r *http.Request) (*User, *User, error) {

	var reqBody User

	if err := util.ParseBody(r, &reqBody); err != nil {
		return &reqBody, nil, err
	}

	if reqBody.Name == "" {
		return &reqBody, nil, fmt.Errorf("invalid user [Name]")
	}

	userData := &User{}

	if err := db.Where("name = ?", reqBody.Name).First(userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &reqBody, nil, fmt.Errorf("user not found")
		}
		return &reqBody, nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &reqBody, userData, nil
}
