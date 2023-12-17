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

type calcReq struct {
	ID                   int    `json:"id"`
	IntervalParam        int    `json:"time_interval"`
	PeoplePerMinuteParam int    `json:"people_per_minute"`
	InputLoad            int    `json:"load"`
	InputToken           string `json:"token"`
}

func Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req := &calcReq{}

	if _, err := json.Marshal(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if req.InputToken != token {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	go worker(req.ID, req.IntervalParam, req.PeoplePerMinuteParam, req.InputLoad)
}

func worker(id, interval, peoplePerMinute, load int) {

	time.Sleep(5 * time.Second)
	calculatedLoad := poissonDistribution(float64(peoplePerMinute)) * interval * load
	fmt.Printf("Simulated Metro Load: %d people\n", calculatedLoad)

	postURL := "your_post_endpoint_url_here" // LOK HERE

	postData := map[string]interface{}{
		"id":             id,
		"calculatedLoad": calculatedLoad,
	}

	// Convert data to JSON
	postBody, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

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
