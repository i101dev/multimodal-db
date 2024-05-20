package routes

import (
	"net/http"

	"github.com/i101dev/multimodal-db/util"

	database "github.com/i101dev/multimodal-db/models/redis"
)

// Alerts

func RegisterAlertRoutes() {

	database.ConnectDB()

	http.HandleFunc("/alerts/all", getAllAlerts)
	http.HandleFunc("/alerts/create", createAlert)
	http.HandleFunc("/alerts/recent", recentAlerts)
}

func getAllAlerts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	allAlerts, err := database.GetAllAlerts(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &allAlerts)
}
func createAlert(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	newAlert, err := database.CreateAlert(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &newAlert)
}
func recentAlerts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	allAlerts, err := database.GetRecentAlerts(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &allAlerts)
}
