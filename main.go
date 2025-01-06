package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CommitResponse struct {
	RequestId string     `json:"request_id"`
	Count     int     `json:"count"`
	Success   int     `json:"success"`
	Fail      int     `json:"fail"`
	Records   []Record `json:"records"`
}

type Record struct {
	Status  string               `json:"status"`
	Index   int                  `json:"index"`
	RecordId string                `json:"record_id"`
	Errors  map[string][]string  `json:"errors"`
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Request received: %s %s\n", r.Method, r.URL)
	fmt.Printf("Request Body: %s\n", string(body))

	responseBody := CommitResponse{
		RequestId: "66666666666666666666666",
		Count:     6666666,
		Success:   6666666,
		Fail:      6666666,
		Records: []Record{
			{
				Status:  "ok",
				Index:   6666666,
				RecordId: "something",
				Errors:  map[string][]string{},
			},
			{
				Status: "error",
				Index:  6666666,
				Errors: map[string][]string{
					"PARENT":   {"Необходимо заполнить «PARENT»."},
					"MERCHANT": {"Значение «XXXXXXXX» для «MERCHANT» уже занято."},
				},
			},
		},
	}

	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/", commitHandler)

	fmt.Println("Listening on http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
