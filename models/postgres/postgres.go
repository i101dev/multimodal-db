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

func NewPostgresConnection() {

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

	init_users(d)

}

func init_users(database *gorm.DB) {

	if err := database.AutoMigrate(&User{}); err != nil {
		fmt.Println("Error initializing [models/users.go]")
		log.Fatal(err)
	}

	db = database
}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

type User struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex" json:"name"`
	UUID   string `json:"uuid"`
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

func CreateUser(r *http.Request) (*User, error) {

	newUser := User{
		UUID:   uuid.New().String(),
		Skills: []Skill{},
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		return nil, err
	}

	result := db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
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
