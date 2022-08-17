package cpfcnpj

import (
	"testing"
	"time"
)

func TestCreateStatus(t *testing.T) {

	tests := []struct {
		give                string
		wantTotalQueryValue uint64
	}{
		{
			give:                "Testa Se Criou Status Como Zero",
			wantTotalQueryValue: 0,
		},
	}

	CreateStatus()
	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			CheckIfEqualInt(t, NumQuery(), tt.wantTotalQueryValue)
		})

	}

}

func TestUpdateStatus(t *testing.T) {

	tests := []struct {
		give                string
		wantTotalQueryValue uint64
	}{
		{
			give:                "Testa Se Incrementou o numero de consultas 1",
			wantTotalQueryValue: 1,
		},
		{
			give:                "Testa Se Incrementou o numero de consultas 2",
			wantTotalQueryValue: 2,
		},
		{
			give:                "Testa Se Incrementou o numero de consultas 3",
			wantTotalQueryValue: 3,
		},
		{
			give:                "Testa Se Incrementou o numero de consultas 4",
			wantTotalQueryValue: 4,
		},
	}

	CreateStatus()
	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			UpdateStatus()
			CheckIfEqualInt(t, NumQuery(), tt.wantTotalQueryValue)
		})

	}

}

func TestUpTimeStatus(t *testing.T) {

	CreateStatus()
	time.Sleep(3 * time.Second)
	lastTime := UptimeQuery()
	CheckIfUptimeIsOK(t, lastTime, 3)
}
