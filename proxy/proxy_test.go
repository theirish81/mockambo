package proxy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/orgs/github/repos", nil)
	out, err := Proxy(req, []string{"http://localhost:8080"}, "https://api.github.com")
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "application/json; charset=utf-8", out.ContentType)
}
