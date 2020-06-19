package entity

import (
	"bytes"
	"encoding/gob"
)

func key(prefix, value string) []byte {
	return []byte(prefix + value)
}

func encode(target interface{}) ([]byte, error) {
	var buff bytes.Buffer

	en := gob.NewEncoder(&buff)

	err := en.Encode(target)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func decode(data []byte, target interface{}) error {
	var buff bytes.Buffer
	de := gob.NewDecoder(&buff)

	_, err := buff.Write(data)
	if err != nil {
		return err
	}

	return de.Decode(target)
}
