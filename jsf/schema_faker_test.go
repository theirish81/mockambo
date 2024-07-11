package jsf

import (
	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"mockambo/extension"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	doc, _ := openapi3.NewLoader().LoadFromFile("../test_data/petstore.yaml")
	path := doc.Paths.Value("/pet/{petId}")
	mext, _ := extension.NewDefaultMextFromExtensions(nil)
	vm := goja.New()
	_ = vm.Set("pathItems", map[string]any{"id": 123})
	out, _ := GenerateDataFromSchema(path.Get.Responses.Value("200").Value.Content.Get("application/json").Schema.Value, mext, vm)
	assert.IsType(t, map[string]any{}, out)
	name, _ := out.(map[string]any)["name"]
	assert.Greater(t, len(name.(string)), 0)
	photoUrls, _ := out.(map[string]any)["photoUrls"]
	assert.IsType(t, []any{}, photoUrls)
}

func TestGenerateString(t *testing.T) {
	var maxLen uint64 = 8
	res, _ := generateString(&openapi3.Schema{
		Type:      &openapi3.Types{"string"},
		MinLength: 3,
		MaxLength: &maxLen,
	}, extension.Mext{})
	assert.GreaterOrEqual(t, len(res), 3)
	assert.LessOrEqual(t, len(res), int(maxLen))

	res, _ = generateString(&openapi3.Schema{
		Type:   &openapi3.Types{"string"},
		Format: "date-time",
	}, extension.Mext{})
	_, err := time.Parse(RFC3339local, res)
	assert.Nil(t, err)
}

func TestGenerateFloat(t *testing.T) {
	res := generateFloat(&openapi3.Schema{
		Type:   &openapi3.Types{"number"},
		Format: "float",
	}, extension.Mext{})
	assert.IsType(t, 1.5, res)

	mn := 0.5
	mx := 1.2
	res = generateFloat(&openapi3.Schema{
		Type:   &openapi3.Types{"number"},
		Format: "float",
		Min:    &mn,
		Max:    &mx,
	}, extension.Mext{})
	assert.GreaterOrEqual(t, res, 0.5)
	assert.LessOrEqual(t, res, 1.2)

	res = generateFloat(&openapi3.Schema{
		Type:   &openapi3.Types{"number"},
		Format: "foo",
	}, extension.Mext{})
	assert.IsType(t, 1.2, res)
}

func TestGenerateInt(t *testing.T) {
	res := generateInt(&openapi3.Schema{
		Type: &openapi3.Types{"integer"},
	}, extension.Mext{})
	assert.IsType(t, 1, res)

	var mn float64 = 1
	var mx float64 = 5
	res = generateInt(&openapi3.Schema{
		Type: &openapi3.Types{"integer"},
		Min:  &mn,
		Max:  &mx,
	}, extension.Mext{})
	assert.GreaterOrEqual(t, res, 1)
	assert.LessOrEqual(t, res, 5)
}
