package yaml

import "github.com/goccy/go-yaml"

type YAML struct{}

func Parser() *YAML {
	return &YAML{}
}

func (y *YAML) Unmarshal(b []byte) (map[string]any, error) {
	var result map[string]any
	if err := yaml.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (y *YAML) Marshal(mp map[string]any) ([]byte, error) {
	return yaml.Marshal(mp)
}
