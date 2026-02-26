package notifications

import (
	"encoding/json"
	"strconv"
	"strings"
)

func unmarshalStringOrNumber(data []byte) (string, bool) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		return "", true
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		return strings.TrimSpace(asString), true
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		return strings.TrimSpace(asNumber.String()), true
	}

	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		return strconv.FormatBool(asBool), true
	}

	return "", false
}

func UnmarshalStringOrNumber(data []byte) string {
	value, ok := unmarshalStringOrNumber(data)
	if !ok {
		return ""
	}

	return strings.TrimSpace(value)
}

func UnmarshalStatusValue(data []byte) string {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		return ""
	}

	if value, ok := unmarshalStringOrNumber(data); ok {
		return strings.TrimSpace(value)
	}

	var asObject struct {
		ID   json.RawMessage `json:"id"`
		Name string          `json:"name"`
	}
	if err := json.Unmarshal(data, &asObject); err != nil {
		return ""
	}

	if name := strings.TrimSpace(asObject.Name); name != "" {
		return name
	}

	return UnmarshalStringOrNumber(asObject.ID)
}
