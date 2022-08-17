package cpfcnpj

import (
	"log"
	"regexp"
)

const (
	SizeToValidTotalCNPJDig1 = 12
	SizeToValidTotalCNPJDig2 = SizeToValidTotalCNPJDig1 + 1

	SizeToValidDig1CNPJ       = 4
	SizeToValidDig2CNPJ       = SizeToValidDig1CNPJ + 1
	SizeToValidDigDefaultCNPJ = SizeToValidDig1CNPJ + SizeToValidDig2CNPJ
	IsCNPJ                    = false
)

func isValidFormatCNPJ(cnpjToCheck string) bool {
	var cnpjRegexp = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}\/?(:?\d{3}[1-9]|\d{2}[1-9]\d|\d[1-9]\d{2}|[1-9]\d{3})-?\d{2}$`)

	if len(cnpjToCheck) < NumTotalDigCPF {
		return false
	}

	return cnpjRegexp.MatchString(cnpjToCheck)
}

//MultiplyNumDigCNPJ os digitos do cnpj por 5 ou 6 *O numero não pode ter caracter especial
func MultiplyNumDigCNPJ(cpfToCheckOnlyNumber string, numIndexFinal int) uint64 {

	multiplicationResult := 0
	strToSum := cpfToCheckOnlyNumber[:numIndexFinal]
	digitMultiplier := numIndexFinal + 1

	for _, nextDigit := range strToSum {
		multiplicationResult += RuneToInt(nextDigit) * digitMultiplier
		digitMultiplier--
	}

	//---Inicio da segunda parte da vaidação do cnpj
	indexLastDigitToCheck := SizeToValidTotalCNPJDig1
	if numIndexFinal == SizeToValidDig2CNPJ {
		indexLastDigitToCheck++
	}

	strCnpjWithoutVerifyDigit := cpfToCheckOnlyNumber[:indexLastDigitToCheck]
	digitMultiplier = SizeToValidDigDefaultCNPJ

	strToSum = strCnpjWithoutVerifyDigit[numIndexFinal:indexLastDigitToCheck]
	for _, nextDigit := range strToSum {
		multiplicationResult += RuneToInt(nextDigit) * digitMultiplier
		digitMultiplier--
	}

	restDivision := multiplicationResult % 11
	compareWithDig := 11 - restDivision

	if restDivision < 2 {
		compareWithDig = 0
	}

	//fmt.Printf("comperToDig2 [%d]FIM\n", compareWithDig)
	return uint64(compareWithDig)
}

func isValidCNPJOnlyValid(cpfToCheck string) bool {

	validDigit1, validDigit2 := VerifyingDigits(cpfToCheck)
	print(validDigit1, validDigit2)

	sumDig1 := MultiplyNumDigCNPJ(cpfToCheck, SizeToValidDig1CNPJ)
	sumDig2 := MultiplyNumDigCNPJ(cpfToCheck, SizeToValidDig2CNPJ)
	print(sumDig1, sumDig2)

	if !ValidateVerifierDigit(sumDig1, validDigit1) {
		log.Printf("Invalid Digit Verifier[%d]\n", validDigit1)
		return false
	}

	return ValidateVerifierDigit(sumDig2, validDigit2)
}

//IsValidCNPJ Check if cnpj is valid
func IsValidCNPJ(cnpjToCheck string) bool {

	if !isValidFormatCNPJ(cnpjToCheck) {
		log.Printf("Invalid Format[%s]\n", cnpjToCheck)
		return false
	}

	cnpjFormated := FormatToValidate(cnpjToCheck)
	return isValidCNPJOnlyValid(cnpjFormated)
}
