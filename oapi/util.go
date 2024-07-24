package oapi

// DeepMerge deeply merges two maps and returns the result.
func DeepMerge(map1, map2 map[any]any) map[any]any {
	result := make(map[any]any)

	for k, v := range map1 {
		result[k] = v
	}

	for k, v2 := range map2 {
		if v1, exists := result[k]; exists {
			result[k] = mergeValues(v1, v2)
		} else {
			result[k] = v2
		}
	}

	return result
}

// mergeValues recursively merges two interface{} values.
func mergeValues(v1 any, v2 any) any {
	switch v1Typed := v1.(type) {
	case map[any]any:
		if v2Typed, ok := v2.(map[any]any); ok {
			return DeepMerge(v1Typed, v2Typed)
		}
	case []any:
		if v2Typed, ok := v2.([]any); ok {
			return append(v1Typed, v2Typed...)
		}
	}
	return v2
}
