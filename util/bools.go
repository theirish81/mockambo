package util

import (
	"math/rand"
)

func RequiredOrRandom(required bool) bool {
	return required || rand.Uint32()%2 == 0
}
