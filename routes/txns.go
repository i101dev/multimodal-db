package routes

import (
	"net/http"

	"github.com/i101dev/multimodal-db/util"

	database "github.com/i101dev/multimodal-db/models/badger"
)

func RegisterTxnRoutes() {

	database.ConnectDB()

	http.HandleFunc("/txn/create", createTxn)
	http.HandleFunc("/txn/getall", getAllTxns)
	http.HandleFunc("/txn/recent", recentTxns)
}

func createTxn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	newTxn, err := database.CreateTxn(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &newTxn)
}
func getAllTxns(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	allTxns, err := database.GetAllTxns(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &allTxns)
}
func recentTxns(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	recentTxns, err := database.GetRecentTxns(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &recentTxns)
}
