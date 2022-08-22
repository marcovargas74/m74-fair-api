package cpfcnpj

import (
	"log"
	"regexp"
)

const (
	NumTotalDigCPF     = 14
	SizeToValidDig1CPF = 9
	SizeToValidDig2CPF = 10
	IsCPF              = true
)

func isValidFormatCPF(cpfToCheck string) bool {
	var CPFRegexp = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)

	if len(cpfToCheck) != NumTotalDigCPF {
		return false
	}

	return CPFRegexp.MatchString(cpfToCheck)
}

//Multiplica os digitos do cpf por 10 ou 11 *O numero não pode ter caracter especial
func MultiplyNumDigCPF(cpfToCheckOnlyNumber string, numIndexFinal int) uint64 {

	strToSum := cpfToCheckOnlyNumber[:numIndexFinal]
	digitMultiplier := (numIndexFinal + 1)

	multiplicationResult := 0
	for _, nextDigit := range strToSum {
		multiplicationResult += RuneToInt(nextDigit) * digitMultiplier
		digitMultiplier--
	}

	restDivision := multiplicationResult % 11
	compareWithDig1 := 11 - restDivision
	if restDivision < 2 {
		compareWithDig1 = 0
	}

	return uint64(compareWithDig1)
}

func isValidCPFOnlyValid(cpfToCheck string) bool {

	validDigit1, validDigit2 := VerifyingDigits(cpfToCheck)

	sumDig1 := MultiplyNumDigCPF(cpfToCheck, SizeToValidDig1CPF)
	sumDig2 := MultiplyNumDigCPF(cpfToCheck, SizeToValidDig2CPF)
	print(sumDig1, sumDig2)

	if !ValidateVerifierDigit(sumDig1, validDigit1) {
		log.Printf("Invalid Digit Verifier[%d]\n", validDigit1)
		return false
	}

	return ValidateVerifierDigit(sumDig2, validDigit2)
}

//IsValidCPF Check if cpf is valid
func IsValidCPF(cpfToCheck string) bool {

	if !isValidFormatCPF(cpfToCheck) {
		log.Printf("Invalid Format[%s]\n", cpfToCheck)
		return false
	}

	cpfFormated := FormatToValidate(cpfToCheck)
	if AllDigitsIsEqual(cpfFormated) {
		log.Printf("Invalid CPF All Digit is Equal[%s]\n", cpfFormated)
		return false
	}

	return isValidCPFOnlyValid(cpfFormated)

}
