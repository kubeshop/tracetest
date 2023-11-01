package functions

import (
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kubeshop/tracetest/server/expression/types"
	"gitlab.com/metakeule/fmtdate"
)

func generateUUID(args ...types.TypedValue) string {
	return gofakeit.UUID()
}

func generateFirstName(args ...types.TypedValue) string {
	return gofakeit.Person().FirstName
}

func generateLastName(args ...types.TypedValue) string {
	return gofakeit.Person().FirstName
}

func generateFullName(args ...types.TypedValue) string {
	person := gofakeit.Person()
	return fmt.Sprintf("%s %s", person.FirstName, person.LastName)
}

func generateEmail(args ...types.TypedValue) string {
	return gofakeit.Email()
}

func generatePhoneNumber(args ...types.TypedValue) string {
	return gofakeit.Phone()
}

func generateCreditCard(args ...types.TypedValue) string {
	return gofakeit.CreditCardNumber(nil)
}

func generateCreditCardCVV(args ...types.TypedValue) string {
	return gofakeit.CreditCardCvv()
}

func generateCreditCardExpiration(args ...types.TypedValue) string {
	return gofakeit.CreditCardExp()
}

func generateRandomInt(args ...types.TypedValue) string {
	min, _ := strconv.Atoi(args[0].Value)
	max, _ := strconv.Atoi(args[1].Value)
	return fmt.Sprintf("%d", gofakeit.Number(min, max))
}

func generateDate(args ...types.TypedValue) string {
	format := time.DateOnly
	if len(args) > 0 {
		format = args[0].Value
	}

	return fmtdate.Format(format, time.Now())
}

func generateDateTime(args ...types.TypedValue) string {
	format := time.RFC3339
	if len(args) > 0 {
		format = args[0].Value
	}
	return fmtdate.Format(format, time.Now())
}
