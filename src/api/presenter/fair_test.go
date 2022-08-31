package presenter

import (
	"net/http"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestSelectKeySearch(t *testing.T) {

	tests := []struct {
		give      string
		inData    string
		wantKey   string
		wantValue string
	}{
		{
			give:      "Fair SelectKeySearch test with empty char",
			inData:    "",
			wantKey:   "",
			wantValue: "",
		},
		{
			give:      "Fair SelectKeySearch test with name key",
			inData:    "/fairs?name=vila",
			wantKey:   "name",
			wantValue: "vila",
		},
		{
			give:      "Fair SelectKeySearch test with district key",
			inData:    "/fairs?district=VILA FORMOSA",
			wantKey:   "district",
			wantValue: "VILA FORMOSA",
		},
		{
			give:      "Fair SelectKeySearch test with region5 key",
			inData:    "/fairs?region5=Sul",
			wantKey:   "region5",
			wantValue: "Sul",
		},
		{
			give:      "Fair SelectKeySearch test with neighborhood key",
			inData:    "/fairs?neighborhood=VILA FORMOSA",
			wantKey:   "neighborhood",
			wantValue: "VILA FORMOSA",
		},
		{
			give:      "Fair SelectKeySearch test with region5= key",
			inData:    "/fairs?region5= ",
			wantKey:   "region5",
			wantValue: " ",
		},
		{
			give:      "Fair SelectKeySearch test with name= key",
			inData:    "/fairs?name= ",
			wantKey:   "name",
			wantValue: " ",
		},
		{
			give:      "Fair SelectKeySearch test with name= fail key",
			inData:    "/fairs?name=",
			wantKey:   "",
			wantValue: "",
		}}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, tt.inData, nil)
			key, value := SelectKeySearch(request)
			assert.Equal(t, key, tt.wantKey)
			assert.Equal(t, value, tt.wantValue)
		})

	}

}
