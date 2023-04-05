package model

import (
	"encoding/json"
)

type Meta struct {
	Size     int    `json:"size"`
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}

const (
	METADATA_SIZE = 1024
)

func (m *Meta) Equals(other *Meta) bool {
	if m.Checksum != other.Checksum ||
		m.Name != other.Name ||
		m.Size != other.Size {
		return false
	}

	return true
}

func (m *Meta) Serialize(width int) ([]byte, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := make([]byte, width)
	for i := range result {
		result[i] = ' '
	}
	copy(result, bs)

	return result, nil
}

func Deserialize(bs []byte) (*Meta, error) {
	meta := &Meta{}

	err := json.Unmarshal(bs, &meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}
