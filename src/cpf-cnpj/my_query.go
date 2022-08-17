package cpfcnpj

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//MyQuery Strutc Main Used in AppValidate
type MyQuery struct {
	ID        string    `json:"id" bson:"id"`
	Number    string    `json:"cpf" bson:"cpf"`
	IsValid   bool      `json:"is_valid" bson:"is_valid"`
	IsCPF     bool      `json:"is_cpf" bson:"is_cpf"`
	IsCNPJ    bool      `json:"is_cnpj" bson:"is_cnpj"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

//ValidCPFQueryGeneric Valid CPF
func (q *MyQuery) ValidCPFQueryGeneric(newCPF string) (int, string) {

	q.Number = newCPF
	q.IsValid = true
	q.IsCNPJ = false

	q.IsCPF = true
	if !IsValidCPF(q.Number) {
		q.IsValid = false
		message := fmt.Sprintf("Something gone wrong: Invalid CPF:%s\n", q.Number)
		log.Println(message)
		return http.StatusBadRequest, message
	}

	return http.StatusOK, "SUCCESS CPF is VALID"
}

//ValidCNPJQueryGeneric Valid CNPJ
func (q *MyQuery) ValidCNPJQueryGeneric(newCNPJ string) (int, string) {

	q.Number = newCNPJ
	q.IsValid = true
	q.IsCPF = false

	q.IsCNPJ = true
	if !IsValidCNPJ(q.Number) {
		q.IsValid = false
		message := fmt.Sprintf("Something gone wrong: Invalid CNPJ:%s\n", q.Number)
		log.Println(message)
		return http.StatusBadRequest, message
	}

	return http.StatusOK, "SUCCESS CNPJ is VALID"
}

//SaveQueryGeneric main fuction to save a new query in system
func (q *MyQuery) SaveQueryGeneric(newCPFofCNPJ string, isCPF bool) (int, string) {

	var code int
	var msg string
	if isCPF {
		code, msg = q.ValidCPFQueryGeneric(newCPFofCNPJ)
	} else {
		code, msg = q.ValidCNPJQueryGeneric(newCPFofCNPJ)
	}

	if code != http.StatusOK {
		return code, msg
	}

	q.ID = NewUUID()
	q.CreatedAt = time.Now()

	if !IsValidUUID(q.ID) {
		message := fmt.Sprintf("Something gone wrong: Invalid ID:%s\n", q.ID)
		log.Println(message)
		return http.StatusBadRequest, message
	}

	result := q.saveQueryInMongoDB()
	if result != nil {
		message := fmt.Sprintf("Can not save cpf/cnpj %v ", q.Number)
		log.Println(message)
		return http.StatusInternalServerError, message
	}

	json, err := q.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	return http.StatusOK, string(json)
}

//QuerysGeneric show All querys save in system
func (q *MyQuery) QuerysGeneric() (int, string) {

	msg, err := q.showQueryAllMongoDB()
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, msg

}

//QuerysByTypeGeneric return ALL CPF or CNPJ pass type in arg
func (q *MyQuery) QuerysByTypeGeneric(isCPF bool) (int, string) {

	msg, err := q.showQuerysByTypeMongoDB(isCPF)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, msg

}

//DeleteQuerysByNum Delete Number
func (q *MyQuery) DeleteQuerysByNumGeneric(findCPForCNPJ string, isCPF bool) (int, string) {

	var code int
	var msg string
	if isCPF {
		code, msg = q.ValidCPFQueryGeneric(findCPForCNPJ)
	} else {
		code, msg = q.ValidCNPJQueryGeneric(findCPForCNPJ)
	}

	if code != http.StatusOK {
		return code, msg
	}

	err := q.deleteQuerysByNumMongoDB(findCPForCNPJ)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, err.Error()
	}

	return http.StatusOK, "SUCCESS TO DELETE CPF/CNPJ"
}

//QuerysByNumGeneric return CPF/CNPJ pass number in arg
func (q *MyQuery) QuerysByNumGeneric(findCPForCNPJ string) (int, string) {

	msg, err := q.showQuerysByNumMongoDB(findCPForCNPJ)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, msg

}
