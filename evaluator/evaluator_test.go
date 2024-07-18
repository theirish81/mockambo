package evaluator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluator_Load(t *testing.T) {
	ev := NewEvaluator()
	data, err := ev.Load("../test_data/sample_payload.json")
	assert.Nil(t, err)
	assert.Greater(t, len(data), 0)
}

func TestEvaluator_RunScript(t *testing.T) {
	ev := NewEvaluator()
	ev.Set("foo", "bar")
	v, err := ev.RunScript("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
}

func TestEvaluator_Template(t *testing.T) {
	ev := NewEvaluator()
	ev.Set("foo", "bar")
	v, err := ev.Template("{{foo}}")
	assert.Nil(t, err)
	assert.Equal(t, "bar", v)
}
