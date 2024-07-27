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

func TestMerger(t *testing.T) {
	merger, _ := NewDoc("../test_data/custom.yaml", "../test_data/merger.yaml")
	_ = merger.Load()
	media := merger.t.Paths.Value("/additional-properties").Get.Responses.Value("200").Value.
		Content.Get("application/json")
	_, ok := media.Extensions["x-mockambo"]
	assert.True(t, ok)
	assert.NotNil(t, media.Example)
}
