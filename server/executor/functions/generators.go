package functions

import (
	"fmt"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
)

func generateUUID(args ...Arg) string {
	return gofakeit.UUID()
}

func generateFirstName(args ...Arg) string {
	return gofakeit.Person().FirstName
}

func generateLastName(args ...Arg) string {
	return gofakeit.Person().FirstName
}

func generateFullName(args ...Arg) string {
	person := gofakeit.Person()
	return fmt.Sprintf("%s %s", person.FirstName, person.LastName)
}

func generateEmail(args ...Arg) string {
	return gofakeit.Email()
}

func generatePhoneNumber(args ...Arg) string {
	return gofakeit.Phone()
}

func generateCreditCard(args ...Arg) string {
	return gofakeit.CreditCardNumber(nil)
}

func generateCreditCardCVV(args ...Arg) string {
	return gofakeit.CreditCardCvv()
}

func generateCreditCardExpiration(args ...Arg) string {
	return gofakeit.CreditCardExp()
}

func generateRandomInt(args ...Arg) string {
	min, _ := strconv.Atoi(args[0].Value)
	max, _ := strconv.Atoi(args[1].Value)
	return fmt.Sprintf("%d", gofakeit.Number(min, max))
}
