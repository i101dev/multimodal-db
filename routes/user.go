package routes

import (
	"net/http"

	"github.com/i101dev/multimodal-db/models/postgres"

	"github.com/i101dev/multimodal-db/util"
)

func RegisterUserRoutes() {

	postgres.ConnectDB()

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
	allUsers, err := postgres.GetAllUsers(r)
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
	newUser, err := postgres.FindUserByID(r)
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
	newUser, err := postgres.CreateUser(r)
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
	newUser, err := postgres.UpdateUser(r)
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
	err := postgres.DeleteUser(r)
	//
	// -----------------------------------------------------------------

	if err != nil {
		util.RespondWithError(w, 500, "Error deleting user")
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("User deleted"))
}
