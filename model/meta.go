package model

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

type Meta struct {
	Size     int    `json:"size"`
	Name     string `json:"name"`
	Checksum []byte `json:"checksum"`
}

const (
	METADATA_SIZE = 1024
)

func GetMetadata(fn string) (*Meta, error) {
	meta := &Meta{}

	file, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file %v", fn)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	meta.Size = int(stat.Size())
	meta.Name = stat.Name()
	meta.Checksum, err = calculateChecksum(file, int(stat.Size()))
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func calculateChecksum(file *os.File, size int) ([]byte, error) {
	hasher := sha256.New()
	buf := make([]byte, 8192)
	reader := bufio.NewReader(file)
	ndone := 0

	for ndone < size {
		nread, err := reader.Read(buf)
		if err != nil {
			return []byte{}, err
		}

		ndone += nread
		nwritten, err := hasher.Write(buf[:nread])
		if err != nil {
			return []byte{}, err
		}
		if nread != nwritten {
			return []byte{}, fmt.Errorf("failed to write all bytes to hasher")
		}
	}

	return hasher.Sum(nil), nil
}

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
