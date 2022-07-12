package utils

import (
	"encoding/json"
	"log"
)

func JSONStatus(msg string) []byte {
	m, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: msg,
	})
	if err != nil {
		log.Println(err)
	}

	return m
}
