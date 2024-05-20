package badger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

// --------------------------------------------------------------------
// --------------------------------------------------------------------

type Txn struct {
	UUID      string `json:"uuid"`
	Item      string `json:"item"`
	Code      string `json:"code"`
	Timestamp int64  `json:"timestamp"`
}

type NullLogger struct{}

func (l *NullLogger) Errorf(string, ...interface{})   {}
func (l *NullLogger) Warningf(string, ...interface{}) {}
func (l *NullLogger) Infof(string, ...interface{})    {}
func (l *NullLogger) Debugf(string, ...interface{})   {}

// --------------------------------------------------------------------
// --------------------------------------------------------------------

const (
	dbPath      = "./tmp/txns"
	genesisData = "First Transaction from Genesis"
)

var (
	badgerDB *badger.DB
)

// --------------------------------------------------------------------
// --------------------------------------------------------------------

func ConnectDB() {

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = &NullLogger{}

	db, err := badger.Open(opts)

	if err != nil {
		log.Fatal("Failed to open BadgerDB:", err)
	}

	badgerDB = db

	fmt.Println("BadgerDB connected successfully")
}

func CreateTxn(r *http.Request) (*Txn, error) {

	var requestBody Txn
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	// -------------------------------------------------------------
	if requestBody.Item == "" {
		return nil, fmt.Errorf("invalid [item]")
	}
	if requestBody.Code == "" {
		return nil, fmt.Errorf("invalid [code]")
	}

	requestBody.UUID = uuid.New().String()
	requestBody.Timestamp = time.Now().Unix()

	// -------------------------------------------------------------
	if err := badgerDB.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(requestBody.UUID), []byte(jsonString(requestBody)))
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	return &requestBody, nil
}

func GetAllTxns(r *http.Request) (*[]Txn, error) {

	var allTxns []Txn

	if err := badgerDB.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("")
		it := txn.NewIterator(opts)
		defer it.Close()

		// -------------------------------------------------------------
		for it.Rewind(); it.Valid(); it.Next() {

			item := it.Item()
			var txnData Txn

			err := item.Value(func(val []byte) error {

				err := json.Unmarshal(val, &txnData)

				if err != nil {
					return fmt.Errorf("failed to deserialize transaction: %v", err)
				}

				allTxns = append(allTxns, txnData)
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
		// -------------------------------------------------------------

	}); err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}

	// -------------------------------------------------------------
	if len(allTxns) == 0 {
		return &allTxns, fmt.Errorf("no transactions yet")
	}

	return &allTxns, nil
}

func GetRecentTxns(r *http.Request) (*[]Txn, error) {

	var requestBody struct {
		Minutes int64 `json:"minutes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	// -------------------------------------------------------------
	var recentTxns []Txn

	currentTime := time.Now().Unix()
	cutoffTime := currentTime - requestBody.Minutes*60

	if err := badgerDB.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		// -------------------------------------------------------------
		for it.Rewind(); it.Valid(); it.Next() {

			item := it.Item()
			var txnData Txn

			err := item.Value(func(val []byte) error {

				err := json.Unmarshal(val, &txnData)

				if err != nil {
					return fmt.Errorf("failed to deserialize transaction: %v", err)
				}

				if txnData.Timestamp >= cutoffTime {
					recentTxns = append(recentTxns, txnData)
				}
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil

		// -------------------------------------------------------------
	}); err != nil {
		return nil, fmt.Errorf("failed to fetch recent transactions: %v", err)
	}

	if len(recentTxns) == 0 {
		return &recentTxns, fmt.Errorf("no recent transactions")
	}

	return &recentTxns, nil
}

func jsonString(data interface{}) string {
	str, _ := json.Marshal(data)
	return string(str)
}
