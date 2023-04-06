package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/cnnrznn/goftp/client"
	"github.com/cnnrznn/goftp/server"
)

func TestE2E(t *testing.T) {
	srcFile := "srcFile.txt"
	dstFile := "dstFile.txt"
	content := []byte("content === content")

	err := os.WriteFile(srcFile, content, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(srcFile)

	server := server.New("localhost:9751", dstFile)
	client := client.New(srcFile, "localhost:9751")
	stopChan := make(chan error)

	go server.Run(stopChan)
	go client.Run(stopChan)

	for err := range stopChan {
		t.Log(err)
	}

	outbs, err := os.ReadFile(dstFile)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(dstFile)

	if !bytes.Equal([]byte(content), outbs) {
		t.Error("files do not match")
	}
}
