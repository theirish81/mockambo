package extension

import (
	"github.com/mitchellh/mapstructure"
	"mockambo/exceptions"
)

const MockamboExt = "x-mockambo"

type Mext struct {
	PayloadGenerationModes []string `yaml:"payloadGenerationModes"`
	Script                 string   `yaml:"script"`
	ValidateRequest        bool     `yaml:"validateRequest"`
	ValidateResponse       bool     `yaml:"validateResponse"`
	Display                bool     `yaml:"display"`
	Proxy                  bool     `yaml:"proxy"`
	ProxyServerIndex       int      `yaml:"ProxyServerIndex"`
	Record                 bool     `yaml:"record"`
	Playback               bool     `yaml:"playback"`
	RecordingKey           string   `yaml:"recordingSignature"`
	RecordingPath          string   `yaml:"recordingPath"`
	LatencyMin             string   `yaml:"latencyMin"`
	LatencyMax             string   `yaml:"latencyMax"`
	ResponseSelector       string   `yaml:"responseSelector"`
	Faker                  string   `yaml:"faker"`
	Template               string   `yaml:"template"`
}

func NewDefaultMextFromExtensions(extensions map[string]any) (Mext, error) {
	mext := Mext{
		PayloadGenerationModes: []string{"script", "template", "faker", "default", "example", "schema"},
		ValidateRequest:        true,
		ValidateResponse:       true,
		Display:                false,
		Proxy:                  false,
		ProxyServerIndex:       0,
		Record:                 false,
		Playback:               false,
		RecordingKey:           "method+'_'+url",
		RecordingPath:          "recording",
		LatencyMin:             "0s",
		LatencyMax:             "0s",
	}
	if ext, ok := extensions[MockamboExt]; ok {
		if err := mapstructure.Decode(ext.(map[string]any), &mext); err != nil {
			return mext, err
		}
	}
	return mext, nil
}

func MergeDefaultMextWithExtensions(def Mext, extensions map[string]any) (Mext, error) {
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
