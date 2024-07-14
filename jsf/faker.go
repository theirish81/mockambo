package jsf

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/dop251/goja"
)

func InstrumentVM(vm *goja.Runtime) {
	_ = vm.Set("fake", Fake)
}

func Fake(t string) any {
	switch t {
	case "address":
		return gofakeit.Address().Address
	case "zip":
		return gofakeit.Zip()
	case "city":
		return gofakeit.City()
	case "country":
		return gofakeit.Country()
	case "countryAbr":
		return gofakeit.CountryAbr()
	case "firstName":
		return gofakeit.FirstName()
	case "lastName":
		return gofakeit.LastName()
	case "domain":
		return gofakeit.DomainName()
	case "url":
		return gofakeit.URL()
	case "email":
		return gofakeit.Email()
	case "creditCard":
		return fmt.Sprintf("%d", gofakeit.CreditCardNumber())
	case "integer":
		return gofakeit.Int32()
	case "float":
		return gofakeit.Float32()
	case "boolean":
		return gofakeit.Bool()
	case "currency":
		return gofakeit.CurrencyLong()
	case "currencyCode":
		return gofakeit.CurrencyShort()
	default:
		return gofakeit.Word()
	}
}
