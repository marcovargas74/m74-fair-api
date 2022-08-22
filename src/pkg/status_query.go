package cpfcnpj

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//StatusQuery VAR used to store number of queries performed since start
var StatusQuery QueryStatus

// //Status used to create a Json Return to endpoint
type QueryStatus struct {
	NumTotalQuery uint64    `json:"num_total_query"`
	StartTime     time.Time `json:"start_time"`
	UpTime        float64   `json:"up_time"`
}

//CreateStatus Create Status
func CreateStatus() {
	StatusQuery.NumTotalQuery = 0
	StatusQuery.StartTime = time.Now()
}

//UpdateStatus increment querys number
func UpdateStatus() {
	StatusQuery.NumTotalQuery++
}

//NumQuery return Total Querys Number
func NumQuery() uint64 {
	return StatusQuery.NumTotalQuery
}

//UptimeQuery return Querys Uptimer
func UptimeQuery() float64 {
	timeElapse := time.Since(StatusQuery.StartTime)

	if timeElapse.Seconds() < 60 {
		return (timeElapse.Seconds())
	}

	fmt.Println("time elapse2:", timeElapse)
	return (timeElapse.Minutes())
}

//ShowStatus Show API Status Qtd Querys and UpTime in segunds
func ShowStatus(w http.ResponseWriter, r *http.Request) {

	StatusQuery.UpTime = UptimeQuery()
	json, err := StatusQuery.MarshalJSON()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
	w.WriteHeader(http.StatusOK)

	fmt.Println("ShowStatus:", string(json))

}

func (q *QueryStatus) MarshalJSON() ([]byte, error) {
	type Alias QueryStatus
	return json.Marshal(&struct {
		*Alias
		StartTime string `json:"start_time"`
	}{
		Alias:     (*Alias)(q),
		StartTime: q.StartTime.Format("02-Jan-06 15:04:05"),
	})
}
