package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// --------------------------------------------------------------------
// --------------------------------------------------------------------

type Alert struct {
	UUID      string `json:"uuid"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

var rdb *redis.Client
var ctx = context.Background()

func ConnectDB() {

	redisHost := os.Getenv("DB_REDIS_HOST")
	redisPort := os.Getenv("DB_REDIS_PORT")
	redisPassword := os.Getenv("DB_REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		log.Fatal("incomplete Redis connection parameters")
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("\n*** >>> Redis connection failed:", err)
		return
	}

	fmt.Println("Redis connected successfully")
}

func CreateAlert(r *http.Request) (*Alert, error) {

	var requestBody Alert

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	// -------------------------------------------------------------
	if requestBody.Title == "" {
		return nil, fmt.Errorf("invalid [title]")
	}
	if requestBody.Body == "" {
		return nil, fmt.Errorf("invalid [body]")
	}

	requestBody.UUID = uuid.New().String()
	requestBody.Timestamp = time.Now().Unix()

	// -------------------------------------------------------------
	alertJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize alert: %v", err)
	}

	err = rdb.Set(ctx, requestBody.UUID, alertJSON, 0).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save alert: %v", err)
	}

	return &requestBody, nil
}

func GetAllAlerts(r *http.Request) (*[]Alert, error) {

	var allAlerts []Alert

	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys: %v", err)
	}

	if len(keys) == 0 {
		return &allAlerts, fmt.Errorf("no alerts yet")
	}

	// -------------------------------------------------------------
	for _, key := range keys {
		alertJSON, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get alert: %v", err)
		}

		var alert Alert
		err = json.Unmarshal([]byte(alertJSON), &alert)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize alert: %v", err)
		}

		allAlerts = append(allAlerts, alert)
	}

	return &allAlerts, nil
}

func GetRecentAlerts(r *http.Request) (*[]Alert, error) {

	var requestBody struct {
		Minutes int64 `json:"minutes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	// -------------------------------------------------------------
	var recentAlerts []Alert

	currentTime := time.Now().Unix()
	cutoffTime := currentTime - requestBody.Minutes*60

	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys: %v", err)
	}

	// -------------------------------------------------------------
	for _, key := range keys {
		alertJSON, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get alert: %v", err)
		}

		var alert Alert
		err = json.Unmarshal([]byte(alertJSON), &alert)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize alert: %v", err)
		}

		if alert.Timestamp >= cutoffTime {
			recentAlerts = append(recentAlerts, alert)
		}
	}

	// -------------------------------------------------------------
	if len(recentAlerts) == 0 {
		return &recentAlerts, fmt.Errorf("no recent alerts")
	}

	return &recentAlerts, nil
}
