package test

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"

	ftp "github.com/cnnrznn/goftp"
)

func TestE2E(t *testing.T) {
	srcFile := "srcFile.txt"
	dstFile := "dstFile.txt"
	content := []byte("content === content\n")

	wg := sync.WaitGroup{}
	wg.Add(1)

	err := os.WriteFile(srcFile, content, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(srcFile)

	go func() {
		err := ftp.ReceiveFile(ftp.Option{
			Addr:     "localhost:9751",
			Filename: dstFile,
		})
		if err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()

	err = ftp.SendFile(ftp.Option{
		Addr:     "localhost:9751",
		Filename: srcFile,
	})
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()

	outbs, err := os.ReadFile(dstFile)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(dstFile)

	if !bytes.Equal([]byte(content), outbs) {
		t.Error("files do not match")
	}
}
