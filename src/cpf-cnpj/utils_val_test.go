package cpfcnpj

import "testing"

func TestFormatToValidateCNPJ(t *testing.T) {

	tests := []struct {
		give        string
		cnpjToCheck string
		wantValue   string
	}{
		{
			give:        "CNPJ format for a string with only digits",
			cnpjToCheck: "000.000.000-00",
			wantValue:   "00000000000",
		},
		{
			give:        "CNPJ format for a string with only digits",
			cnpjToCheck: "00.000.000/0000-00",
			wantValue:   "00000000000000",
		},

		{
			give:        "CNPJ format for a string with only digits",
			cnpjToCheck: "36.562.098/0001-18",
			wantValue:   "36562098000118",
		},
		{
			give:        "CNPJ format for a string with only digits",
			cnpjToCheck: "12.074.074/0001-51",
			wantValue:   "12074074000151",
		},
		{
			give:        "CNPJ format for a string with only digits",
			cnpjToCheck: "24.572.400/0001-30",
			wantValue:   "24572400000130",
		},
		{
			give:        "CNPJ format for a string with only digits",
			cnpjToCheck: "47.425.683/0001-92",
			wantValue:   "47425683000192",
		},
	}
	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			result := FormatToValidate(tt.cnpjToCheck)
			CheckIfEqualString(t, result, tt.wantValue)
		})

	}

}

func TestFormatToValidateCPF(t *testing.T) {

	tests := []struct {
		give       string
		cpfToCheck string
		wantValue  string
	}{
		{
			give:       "CPF format for a string with only digits",
			wantValue:  "00000000000",
			cpfToCheck: "000.000.000-00",
		},
		{
			give:       "CPF format for a string with only digits",
			cpfToCheck: "111.111.111-11",
			wantValue:  "11111111111",
		},

		{
			give:       "CPF format for a string with only digits",
			cpfToCheck: "838.461.722-86",
			wantValue:  "83846172286",
		},
		{
			give:       "CPF format for a string with only digits",
			cpfToCheck: "313.396.023-77",
			wantValue:  "31339602377",
		},
		{
			give:       "CPF format for a string with only digits",
			cpfToCheck: "682.511.941-99",
			wantValue:  "68251194199",
		},
	}

	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			result := FormatToValidate(tt.cpfToCheck)
			CheckIfEqualString(t, result, tt.wantValue)
		})

	}

}

func TestValidateVerifierDigit(t *testing.T) {

	tests := []struct {
		give      string
		wantValue bool

		sumDig1     uint64
		validDigit1 uint64
	}{
		{
			give:        "Test check is number is de same",
			wantValue:   false,
			sumDig1:     0,
			validDigit1: 1,
		},
		{
			give:        "Test check is number is de same",
			wantValue:   false,
			sumDig1:     9,
			validDigit1: 8,
		},
		{
			give:        "Test check is number is de same",
			wantValue:   true,
			sumDig1:     7,
			validDigit1: 7,
		},
		{
			give:        "Test check is number is de same",
			wantValue:   true,
			sumDig1:     0,
			validDigit1: 0,
		},
	}

	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {

			result := ValidateVerifierDigit(tt.sumDig1, tt.validDigit1)
			CheckIfEqualBool(t, result, tt.wantValue)

		})

	}

}

func TestAllDigitsIsEqual(t *testing.T) {

	tests := []struct {
		give       string
		wantValue  bool
		cpfToCheck string
	}{
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "b1080263",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "83846172286",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "31339602377",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "68251194199",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "28875224340",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  false,
			cpfToCheck: "48416241201",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  true,
			cpfToCheck: "00000000000",
		},
		{
			give:       "Test if all digits is Equal - This is a Invalid CPF",
			wantValue:  true,
			cpfToCheck: "11111111111",
		},
	}

	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			result := AllDigitsIsEqual(tt.cpfToCheck)
			CheckIfEqualBool(t, result, tt.wantValue)
		})

	}

}

func TestGetVerifyingDigits(t *testing.T) {

	tests := []struct {
		give       string
		wantValue1 uint64
		wantValue2 uint64
		cpfToCheck string
	}{
		{
			give:       "Get Digits To check if arg is Zeros Numbers",
			wantValue1: 0,
			wantValue2: 0,
			cpfToCheck: "000.000.000-00",
		},
		{
			give:       "Get Digits To check ",
			wantValue1: 1,
			wantValue2: 1,
			cpfToCheck: "111.111.111-11",
		},

		{
			give:       "Get Digits To check ",
			wantValue1: 8,
			wantValue2: 6,
			cpfToCheck: "838.461.722-86",
		},
		{
			give:       "Get Digits To check ",
			wantValue1: 7,
			wantValue2: 7,
			cpfToCheck: "313.396.023-77",
		},
		{
			give:       "Get Digits To check ",
			wantValue1: 9,
			wantValue2: 9,
			cpfToCheck: "682.511.941-99",
		},
	}

	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			Dig1, Dig2 := VerifyingDigits(tt.cpfToCheck)
			CheckIfEqualInt(t, Dig1, tt.wantValue1)
			CheckIfEqualInt(t, Dig2, tt.wantValue2)
		})

	}

}

func TestVerifyingDigitsCNPJ(t *testing.T) {

	tests := []struct {
		give        string
		wantValue1  uint64
		wantValue2  uint64
		cnpjToCheck string
	}{
		{
			give:        "Get Digits To check if arg is Zeros Numbers",
			wantValue1:  0,
			wantValue2:  0,
			cnpjToCheck: "00.000.000/0000-00",
		},
		{
			give:        "Get Digits To check ",
			wantValue1:  1,
			wantValue2:  8,
			cnpjToCheck: "36.562.098/0001-18",
		},

		{
			give:        "Get Digits To check ",
			wantValue1:  5,
			wantValue2:  1,
			cnpjToCheck: "12.074.074/0001-51",
		},
		{
			give:        "Get Digits To check ",
			wantValue1:  3,
			wantValue2:  0,
			cnpjToCheck: "24.572.400/0001-30",
		},
		{
			give:        "Get Digits To check ",
			wantValue1:  9,
			wantValue2:  2,
			cnpjToCheck: "47.425.683/0001-92",
		},
	}

	for _, tt := range tests {

		t.Run(tt.give, func(t *testing.T) {
			Dig1, Dig2 := VerifyingDigits(tt.cnpjToCheck)
			CheckIfEqualInt(t, Dig1, tt.wantValue1)
			CheckIfEqualInt(t, Dig2, tt.wantValue2)
		})

	}

}
