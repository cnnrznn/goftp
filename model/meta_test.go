package model

import (
	"testing"
)

func TestSerialize(t *testing.T) {
	model := &Meta{
		Name:     "myfile.txt",
		Size:     4000000,
		Checksum: "bullshit",
	}

	bs, err := model.Serialize(METADATA_SIZE)
	if err != nil {
		t.Error(err)
	}

	if len(bs) != METADATA_SIZE {
		t.Errorf("serialized payload is the wrong size: got %v, want %v", len(bs), METADATA_SIZE)
	}

	meta, err := Deserialize(bs)
	if err != nil {
		t.Error(err)
	}

	if !model.Equals(meta) {
		t.Errorf("objects do not match: %v\n%v\n", model, meta)
	}
}
