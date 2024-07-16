package jsf

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestFake(t *testing.T) {
	assert.IsType(t, int32(32), Fake("integer"))
	assert.IsType(t, float32(32.5), Fake("float"))
	assert.IsType(t, "foobar", Fake("address"))
	assert.IsType(t, "foobar", Fake("sblargh"))
	assert.IsType(t, true, Fake("boolean"))

	_, err := url.Parse(Fake("url").(string))
	assert.Nil(t, err)
}
