package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {

	// if code > 499 {
	// 	log.Println("Responding with 5XX error:", msg)
	// }

	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Printf("failed to encode order to JSON: %+v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func ParseJSONRequestBody(r *http.Request) (map[string]interface{}, error) {

	var requestBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return requestBody, nil
}

func ParseBody(r *http.Request, x interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(x); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}
	return nil
}
