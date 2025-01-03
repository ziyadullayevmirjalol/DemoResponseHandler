package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CommitResponse struct {
	RequestId int     `json:"request_id"`
	Count     int     `json:"count"`
	Success   int     `json:"success"`
	Fail      int     `json:"fail"`
	Records   []Record `json:"records"`
}

type Record struct {
	Status  string               `json:"status"`
	Index   int                  `json:"index"`
	RecordId *int                `json:"record_id,omitempty"`
	Errors  map[string][]string  `json:"errors"`
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Log the request body
	fmt.Printf("Request received: %s %s\n", r.Method, r.URL)
	fmt.Printf("Request Body: %s\n", string(body))

	// Prepare the response body
	responseBody := CommitResponse{
		RequestId: 999,
		Count:     2,
		Success:   0,
		Fail:      2,
		Records: []Record{
			{
				Status:  "ok",
				Index:   0,
				RecordId: &[]int{999}[0],
				Errors:  map[string][]string{},
			},
			{
				Status: "error",
				Index:  1,
				Errors: map[string][]string{
					"PARENT":   {"Необходимо заполнить «PARENT»."},
					"MERCHANT": {"Значение «XXXXXXXX» для «MERCHANT» уже занято."},
				},
			},
		},
	}

	// Serialize response to JSON
	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	// Set response headers and content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the response
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/", commitHandler)

	// Start the server
	fmt.Println("Listening on http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
