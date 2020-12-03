package jsonf

import (
	"bytes"
	"encoding/json"
)

func Encode(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Decode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

func DecodeUseNumber(data []byte, val interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(&val)
}
