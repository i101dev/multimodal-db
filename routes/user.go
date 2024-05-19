package routes

import (
	"net/http"

	// database "github.com/i101dev/multimodal-db/models/postgres"
	database "github.com/i101dev/multimodal-db/models/mysql"

	"github.com/i101dev/multimodal-db/util"
)

func RegisterUserRoutes() {

	database.ConnectDB()

	http.HandleFunc("/users", getAll)
	http.HandleFunc("/users/find", find)
	http.HandleFunc("/users/create", create)
	http.HandleFunc("/users/update", update)
	http.HandleFunc("/users/delete", delete)
}
func getAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	allUsers, err := database.GetAllUsers(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &allUsers)
}

func find(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	newUser, err := database.FindUserByID(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, "Error finding user")
		return
	}

	util.RespondWithJSON(w, 200, &newUser)
}

func create(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	newUser, err := database.CreateUser(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, &newUser)
}

func update(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	newUser, err := database.UpdateUser(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, "Error updating user")
		return
	}

	util.RespondWithJSON(w, 200, &newUser)
}

func delete(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// -----------------------------------------------------------------
	//
	err := database.DeleteUser(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, "Error deleting user")
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("User deleted"))
}
