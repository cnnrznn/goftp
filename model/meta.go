package model

import (
	"encoding/json"
	"fmt"
)

type Meta struct {
	Size     int    `json:"size"`
	Name     string `json:"name"`
	Checksum []byte `json:"checksum"`
}

const (
	METADATA_SIZE = 1024
)

func (m *Meta) ChecksumEquals(other []byte) bool {
	if len(m.Checksum) != len(other) {
		return false
	}

	for i := 0; i < len(m.Checksum); i++ {
		if m.Checksum[i] != other[i] {
			return false
		}
	}

	return true
}

func (m *Meta) String() string {
	return fmt.Sprintf("%d, %s", m.Size, m.Name)
}

func (m *Meta) Equals(other *Meta) bool {
	if !m.ChecksumEquals(other.Checksum) ||
		m.Name != other.Name ||
		m.Size != other.Size {
		return false
	}

	return true
}

func Serialize(m *Meta) ([]byte, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := make([]byte, METADATA_SIZE)
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
