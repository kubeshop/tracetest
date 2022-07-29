package functions

import (
	"fmt"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
)

func generateUUID(args ...FunctionArg) string {
	return gofakeit.UUID()
}

func generateFirstName(args ...FunctionArg) string {
	return gofakeit.Person().FirstName
}

func generateLastName(args ...FunctionArg) string {
	return gofakeit.Person().FirstName
}

func generateFullName(args ...FunctionArg) string {
	person := gofakeit.Person()
	return fmt.Sprintf("%s %s", person.FirstName, person.LastName)
}

func generateEmail(args ...FunctionArg) string {
	return gofakeit.Email()
}

func generatePhoneNumber(args ...FunctionArg) string {
	return gofakeit.Phone()
}

func generateCreditCard(args ...FunctionArg) string {
	return gofakeit.CreditCardNumber(nil)
}

func generateCreditCardCVV(args ...FunctionArg) string {
	return gofakeit.CreditCardCvv()
}

func generateCreditCardExpiration(args ...FunctionArg) string {
	return gofakeit.CreditCardExp()
}

func generateRandomInt(args ...FunctionArg) string {
	min, _ := strconv.Atoi(args[0].Value)
	max, _ := strconv.Atoi(args[1].Value)
	return fmt.Sprintf("%d", gofakeit.Number(min, max))
}
