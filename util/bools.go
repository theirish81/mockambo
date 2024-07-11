package util

import "github.com/brianvoe/gofakeit"

func RequiredOrRandom(required bool) bool {
	return required || gofakeit.Bool()

}
