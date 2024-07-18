package extension

import (
	"github.com/mitchellh/mapstructure"
	"mockambo/exceptions"
)

const MockamboExt = "x-mockambo"

// Mext is the data structure representing the x-mockambo OpenAPI extension.
// No fields are pointers to simplify the effects of cloning
type Mext struct {
	PayloadGenerationModes   []string `yaml:"payloadGenerationModes"`
	Script                   string   `yaml:"script"`
	ValidateRequest          bool     `yaml:"validateRequest"`
	ValidateResponse         bool     `yaml:"validateResponse"`
	Display                  bool     `yaml:"display"`
	Proxy                    bool     `yaml:"proxy"`
	ProxyServerIndex         int      `yaml:"proxyServerIndex"`
	Record                   bool     `yaml:"record"`
	Playback                 bool     `yaml:"playback"`
	RecordingSignatureScript string   `yaml:"recordingSignatureScript"`
	RecordingPath            string   `yaml:"recordingPath"`
	LatencyMin               string   `yaml:"latencyMin"`
	LatencyMax               string   `yaml:"latencyMax"`
	ResponseSelector         string   `yaml:"responseSelector"`
	Faker                    string   `yaml:"faker"`
	Template                 string   `yaml:"template"`
}

// NewMextFromExtensions will create a default mext and merge it with the root x-mockambo extension (if present).
// Mind that the extensions argument is the entire extensions map, the method will look for the x-mockambo extension
// within the extensions map
func NewMextFromExtensions(extensions map[string]any) (Mext, error) {
	mext := Mext{
		PayloadGenerationModes:   []string{"script", "template", "faker", "default", "example", "schema"},
		ValidateRequest:          true,
		ValidateResponse:         true,
		Display:                  false,
		Proxy:                    false,
		ProxyServerIndex:         0,
		Record:                   false,
		Playback:                 false,
		RecordingSignatureScript: "method+'_'+url",
		RecordingPath:            "recording",
		LatencyMin:               "0s",
		LatencyMax:               "0s",
	}
	if ext, ok := extensions[MockamboExt]; ok {
		if err := mapstructure.Decode(ext.(map[string]any), &mext); err != nil {
			return mext, err
		}
	}
	return mext, nil
}

// MergeMextWithExtensions merges an existing Mext with a new extensions map
// Mind that the extensions argument is the entire extensions map, the method will look for the x-mockambo extension
// within ig
func MergeMextWithExtensions(def Mext, extensions map[string]any) (Mext, error) {
	def.Display = false
	modes := def.PayloadGenerationModes
	def.PayloadGenerationModes = make([]string, len(def.PayloadGenerationModes))
	copy(def.PayloadGenerationModes, modes)
	if ext, ok := extensions[MockamboExt]; ok {
		if err := mapstructure.Decode(ext.(map[string]any), &def); err != nil {
			return def, exceptions.Wrap("decode_extension", err)
		}
	}
	return def, nil
}
