package db

import (
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestUpsert(t *testing.T) {
	px := path.Join(os.TempDir(), gofakeit.UUID())
	_ = Upsert("foo", []byte("bar"), px)
	res, _ := Get("foo", px)
	assert.Equal(t, []byte("bar"), res)
	res, err := Get("DJANGA", px)
	assert.NotNil(t, err)
}
