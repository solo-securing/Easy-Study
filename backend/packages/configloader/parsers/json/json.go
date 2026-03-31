package json

import "encoding/json"

type JSON struct{}

func Parser() *JSON {
	return &JSON{}
}

func (j *JSON) Unmarshal(b []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (j *JSON) Marshal(mp map[string]any) ([]byte, error) {
	return json.Marshal(mp)
}
