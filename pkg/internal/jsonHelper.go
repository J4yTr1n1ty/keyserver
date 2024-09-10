package internal

import (
	"encoding/json"
)

func GetJsonString(data any) (string, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func GetJsonBytes(data any) ([]byte, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json, nil
}
