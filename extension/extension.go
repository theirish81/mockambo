package extension

import "github.com/mitchellh/mapstructure"

const MockamboExt = "x-mockambo"

type Mext struct {
	PayloadGenerationModes []string `yaml:"payloadGenerationModes"`
	Script                 *string  `yaml:"script"`
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
	ResponseSelector       *string  `yaml:"responseSelector"`
	Faker                  *string  `yaml:"faker"`
}

func NewDefaultMextFromExtensions(extensions map[string]any) (Mext, error) {
	mext := Mext{
		PayloadGenerationModes: []string{"script", "default", "example", "schema"},
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
	def.Script = nil
	def.Display = false
	def.ResponseSelector = nil
	def.Faker = nil
	modes := def.PayloadGenerationModes
	def.PayloadGenerationModes = make([]string, len(def.PayloadGenerationModes))
	copy(def.PayloadGenerationModes, modes)
	if ext, ok := extensions[MockamboExt]; ok {
		if err := mapstructure.Decode(ext.(map[string]any), &def); err != nil {
			return def, err
		}
	}
	return def, nil
}
