package utils

import "encoding/json"

func JSONStatus(msg string) ([]byte, error) {
	m, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: msg,
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}
