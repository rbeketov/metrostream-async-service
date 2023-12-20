package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var token = "Hg12HdEdEiid9-djEDegE"

type modelParam struct {
	ModelID   int `json:"model_id"`
	InputLoad int `json:"load"`
}

type outputParam struct {
	ModelID    int `json:"model_id"`
	OutputLoad int `json:"output_load"`
}

type calcReq struct {
	ID                   int          `json:"id"`
	IntervalParam        int          `json:"time_interval"`
	PeoplePerMinuteParam int          `json:"people_per_minute"`
	Params               []modelParam `json:"modelings"`
	InputToken           string       `json:"token"`
}

func Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req := &calcReq{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.InputToken != token {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	go worker(req.ID, req.IntervalParam, req.PeoplePerMinuteParam, req.Params)

}

func worker(id int, interval int, peoplePerMinute int, params []modelParam) {
	time.Sleep(6 * time.Second)

	var results []outputParam

	for _, param := range params {
		resultLoading := poissonDistribution(float64(peoplePerMinute)) * interval * param.InputLoad / 100

		// Создаем новый элемент outputParam и добавляем в слайс results
		result := outputParam{
			ModelID:    param.ModelID,
			OutputLoad: resultLoading,
		}
		results = append(results, result)
	}
	// fmt.Printf("Simulated Metro Load for %s: %d people\n", param.ModelID, resultLoading)

	// Create the JSON payload
	requestData := map[string]interface{}{
		"application_id": id,
		"results":        results,
		"token":          token,
	}

	// Convert data to JSON
	postBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	postURL := "http://localhost:8000/api/applications/write_result_modeling/"

	// Make the HTTP POST request
	resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Unexpected status code %d\n", resp.StatusCode)
		return
	}
}

func poissonDistribution(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}

	return k - 1
}
