package routes

import (
	"net/http"

	"github.com/i101dev/multimodal-db/util"
)

// ------------------------------------------------------------------------
// Routes -----------------------------------------------------------------

func RegisterTestRoutes() {
	http.HandleFunc("/testGet", testGet)
	http.HandleFunc("/testPut", testPut)
	http.HandleFunc("/testPost", testPost)
}

// ------------------------------------------------------------------------
// Handlers ---------------------------------------------------------------

func testGet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	name := params.Get("name")

	if name == "" {
		name = "World"
	}

	body := map[string]interface{}{
		"message": "Hello, " + name + "! This is a GET request.",
	}

	util.RespondWithJSON(w, 200, body)
}

func testPost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestBody, err := util.ParseJSONRequestBody(r)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	body := map[string]interface{}{
		"message": "Hello, World! This is a POST request.",
		"data":    requestBody,
	}

	util.RespondWithJSON(w, 200, body)
}

func testPut(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		util.RespondWithError(w, 400, "Method not allowed")
		return
	}

	params := r.URL.Query()

	name := params.Get("name")
	if name == "" {
		name = "World"
	}

	body := map[string]interface{}{
		"message": "Hellow, " + name + "! This is a PUT request.",
	}

	util.RespondWithJSON(w, 200, body)
}
