package util

import (
	"math/rand"
	"mockambo/exceptions"
	"mockambo/extension"
	"time"
)

func ComputeLatency(mext extension.Mext, req *Request) (time.Duration, error) {
	sleepTime := 0 * time.Second
	elapsed := time.Now().Sub(req.CreatedAt)
	latencyMin, err := time.ParseDuration(mext.LatencyMin)
	if err != nil {
		return sleepTime, exceptions.Wrap("parse_latency_min", err)
	}
	latencyMax, err := time.ParseDuration(mext.LatencyMax)
	if err != nil {
		return sleepTime, exceptions.Wrap("parse_latency_max", err)
	}
	if elapsed < latencyMin {
		rx := latencyMax - latencyMin
		delay := int(latencyMin.Milliseconds()) + rand.Intn(int(rx.Milliseconds()))
		sleepTime = time.Duration(delay) * time.Millisecond
		elapsed = sleepTime
	} else if elapsed < latencyMax {
		rx := latencyMax - elapsed
		sleepTime = time.Duration(rand.Intn(int(rx.Milliseconds()))) * time.Millisecond
	}
	return sleepTime, nil
}
