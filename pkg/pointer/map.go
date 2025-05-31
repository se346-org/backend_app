package pointer

import "encoding/json"

func ToMap[T any](data T) (map[string]any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
