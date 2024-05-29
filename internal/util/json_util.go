package util

import (
	"encoding/json"
	"log"
)

func ToJson(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatal("Failed to marshal...")
	}
	return string(bytes)
}

func ToPrettyJson(v any) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal...")
	}
	return string(bytes)
}
