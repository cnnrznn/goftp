package model

import (
	"testing"
)

func TestSerialize(t *testing.T) {
	meta := &Meta{
		Name:     "myfile.txt",
		Size:     4000000,
		Checksum: []byte("bullshit"),
	}

	bs, err := Serialize(meta)
	if err != nil {
		t.Error(err)
	}

	if len(bs) != METADATA_SIZE {
		t.Errorf("serialized payload is the wrong size: got %v, want %v", len(bs), METADATA_SIZE)
	}

	metaResult, err := Deserialize(bs)
	if err != nil {
		t.Error(err)
	}

	if !meta.Equals(metaResult) {
		t.Errorf("objects do not match: %v\n%v\n", meta, metaResult)
	}
}
