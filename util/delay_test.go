package util

import (
	"github.com/stretchr/testify/assert"
	"mockambo/extension"
	"net/http"
	"testing"
	"time"
)

func TestComputeLatency(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://www.example.com", nil)
	req := NewRequest(r)
	duration, _ := ComputeLatency(extension.Mext{
		LatencyMin: "500ms",
		LatencyMax: "1s",
	}, req)
	assert.Greater(t, duration, 500*time.Millisecond)
	assert.Less(t, duration, 1*time.Second)
}
