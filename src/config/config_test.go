package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestConfigGetMysqlURL(t *testing.T) {

	tests := []struct {
		give         string
		wantUser     string
		wantPass     string
		wantDatabase string
		wantAddr     string
		wantPort     string
	}{
		{
			give:         "Config MySQL URL Default VALUES",
			wantUser:     DEFAULT_DB_USER,
			wantPass:     DEFAULT_DB_PASSWORD,
			wantDatabase: DEFAULT_DB_DATABASE,
			wantAddr:     DEFAULT_DB_ADDRESS,
			wantPort:     DEFAULT_DB_PORT,
		},
		{
			give:         "Config MySQL URL NOT Default VALUES",
			wantUser:     "UserTest",
			wantPass:     "PasswordTest",
			wantDatabase: "DatabaseTest",
			wantAddr:     "127.0.0.1",
			wantPort:     "8800",
		},
		{
			give:         "Config MySQL URL WITH spaces VALUES",
			wantUser:     " ",
			wantPass:     " ",
			wantDatabase: " ",
			wantAddr:     " ",
			wantPort:     " ",
		}}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			os.Setenv(_DB_USER, tt.wantUser)
			os.Setenv(_DB_PASSWORD, tt.wantPass)
			os.Setenv(_DB_DATABASE, tt.wantDatabase)
			os.Setenv(_DB_ADDRESS, tt.wantAddr)
			os.Setenv(_DB_PORT, tt.wantPort)

			wantURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
				tt.wantUser, tt.wantPass, tt.wantAddr, tt.wantPort, tt.wantDatabase)

			url, _ := ConfigGetMysqlURL()
			assert.Equal(t, url, wantURL)
		})

		os.Unsetenv(_DB_USER)
		os.Unsetenv(_DB_PASSWORD)
		os.Unsetenv(_DB_DATABASE)
		os.Unsetenv(_DB_ADDRESS)
		os.Unsetenv(_DB_PORT)

	}

}

func TestConfigGetMysqlURLDefault(t *testing.T) {
	wantURL := DEFAULT_URL_MYSQL

	os.Clearenv()
	url, _ := ConfigGetMysqlURL()
	assert.Equal(t, url, wantURL)
}

func TestConfigGetAPIGeneral(t *testing.T) {

	tests := []struct {
		give        string
		wantTypeApp string
		wantPortMem string
		wantPortSQL string
		wantLogFile string
		wantIsProd  bool
	}{
		{
			give:        "Config General Default VALUES",
			wantTypeApp: TYPE_PROD,
			wantPortMem: ":5000",
			wantPortSQL: ":5001",
			wantLogFile: "./fairAPI.log",
			wantIsProd:  true,
		},
		{
			give:        "Config General DEV VALUES",
			wantTypeApp: TYPE_DEV,
			wantPortMem: ":6000",
			wantPortSQL: ":6001",
			wantLogFile: " ",
			wantIsProd:  false,
		},
		{
			give:        "Config General TESTs VALUES",
			wantTypeApp: TYPE_TEST,
			wantPortMem: ":0000",
			wantPortSQL: ":0001",
			wantLogFile: "./test.txt",
			wantIsProd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			os.Setenv(_TYPE_APP, tt.wantTypeApp)
			os.Setenv(_SERVER_API_PORT_MEM, tt.wantPortMem)
			os.Setenv(_SERVER_API_PORT_SQL, tt.wantPortSQL)
			os.Setenv(_LOG_FILE, tt.wantLogFile)

			tc, _ := ConfigGetAPIGeneral()
			assert.Equal(t, tc.APITypeApp, tt.wantTypeApp)
			assert.Equal(t, tc.APIServerPortMem, tt.wantPortMem)
			assert.Equal(t, tc.APIServerPortSQL, tt.wantPortSQL)
			assert.Equal(t, tc.APILogFile, tt.wantLogFile)
			assert.Equal(t, tc.IsProdType(), tt.wantIsProd)
		})
		os.Unsetenv(_TYPE_APP)
		os.Unsetenv(_SERVER_API_PORT_MEM)
		os.Unsetenv(_SERVER_API_PORT_SQL)
		os.Unsetenv(_LOG_FILE)
	}

}

func TestConfigAPIGeneralDefault(t *testing.T) {
	var want ConfigAPI
	want.APITypeApp = TYPE_PROD
	want.APIServerPortMem = DEFAULT_SERVER_API_PORT_MEM
	want.APIServerPortSQL = DEFAULT_SERVER_API_PORT_SQL
	want.APILogFile = DEFAULT_LOG_FILE

	os.Clearenv()
	tc, _ := ConfigGetAPIGeneral()
	assert.Equal(t, tc.APITypeApp, want.APITypeApp)
	assert.Equal(t, tc.APIServerPortMem, want.APIServerPortMem)
	assert.Equal(t, tc.APIServerPortSQL, want.APIServerPortSQL)
	assert.Equal(t, tc.APILogFile, want.APILogFile)

	assert.Equal(t, tc.IsProdType(), true)

}
