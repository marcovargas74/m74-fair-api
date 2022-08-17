package cpfcnpj

import (
	"testing"
	"time"
)

const (
	erroMsg = "Test fail got Value[%v], wait Value [%v]"
)

//CheckIfEqualString check if result is OK type string
func CheckIfEqualString(t *testing.T, gotValue, waitValue string) {
	t.Helper()
	if gotValue != waitValue {
		t.Errorf(erroMsg, gotValue, waitValue)
	}
}

//CheckIfEqualBool check if result is OK type BOOL
func CheckIfEqualBool(t *testing.T, gotValue, waitValue bool) {
	t.Helper()
	if gotValue != waitValue {
		t.Errorf(erroMsg, gotValue, waitValue)
	}
}

//CheckIfEqualInt check if result is OK type INT
func CheckIfEqualInt(t *testing.T, gotValue, waitValue uint64) {
	t.Helper()
	if gotValue != waitValue {
		t.Errorf(erroMsg, gotValue, waitValue)
	}
}

//CheckIfEqualFloat check if result is OK type FLOAT
func CheckIfEqualFloat(t *testing.T, gotValue, waitValue float64) {
	t.Helper()
	if gotValue != waitValue {
		t.Errorf(erroMsg, gotValue, waitValue)
	}
}

//CheckIfEqualTime check if result is OK type Time
func CheckIfEqualTime(t *testing.T, gotValue, waitValue time.Time) {
	t.Helper()
	if gotValue != waitValue {
		t.Errorf(erroMsg, gotValue, waitValue)
	}
}

//CheckIfUptimeIsOK check if result is OK to UpTime
func CheckIfUptimeIsOK(t *testing.T, gotValue, waitValue float64) {
	t.Helper()
	if gotValue <= waitValue {
		t.Errorf(erroMsg, gotValue, waitValue)
	}
}

func TestSetUsingMongoDocker(t *testing.T) {

	tests := []struct {
		give      string
		wantValue bool
		setFlag   bool
	}{
		{
			give:      "Test SetOpen Flag end Get Var IsUsingMongoDB",
			wantValue: false,
			setFlag:   false,
		},
		{
			give:      "Test SetOpen Flag end Get Var IsUsingMongoDB",
			wantValue: true,
			setFlag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			SetUsingMongoDocker(tt.setFlag)
			result := CheckIsUsingMongoDocker()
			CheckIfEqualBool(t, result, tt.wantValue)
		})

	}

}
