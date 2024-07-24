package oapi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeepMerge(t *testing.T) {
	m1 := map[any]any{
		"foo": map[any]any{
			"bar": 1,
			"tar": 2,
		},
		"arr": []any{
			1,
			2,
			3,
		},
	}
	m2 := map[any]any{
		"foo": map[any]any{
			"zar": 3,
		},
		"test": "abc",
		"arr": []any{
			4,
		},
	}

	out := DeepMerge(m1, m2)
	assert.Equal(t, 1, out["foo"].(map[any]any)["bar"])
	assert.Equal(t, 3, out["foo"].(map[any]any)["zar"])
	assert.Equal(t, "abc", out["test"])
	assert.Equal(t, 1, out["arr"].([]any)[0])
	assert.Equal(t, 2, out["arr"].([]any)[1])
	assert.Equal(t, 3, out["arr"].([]any)[2])
	assert.Equal(t, 4, out["arr"].([]any)[3])

}
