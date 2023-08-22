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

	err := os.WriteFile(srcFile, content, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(srcFile)
	defer os.Remove(dstFile)

	wg := sync.WaitGroup{}
	wg.Add(1)

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
		Retries:  3,
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

func TestRepeated(t *testing.T) {
	srcFile := "srcFile.txt"
	dstFile := "dstFile.txt"
	content := []byte("content === content\n")
	iter := 10

	err := os.WriteFile(srcFile, content, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(srcFile)
	defer os.Remove(dstFile)

	wg := sync.WaitGroup{}
	wg.Add(iter)

	go func() {
		for i := 0; i < iter; i++ {
			fmt.Printf("%v\r", i)
			err := ftp.ReceiveFile(ftp.Option{
				Addr:     "localhost:9752",
				Filename: dstFile,
			})
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}
	}()

	for i := 0; i < iter; i++ {
		err = ftp.SendFile(ftp.Option{
			Addr:     "localhost:9752",
			Filename: srcFile,
			Retries:  3,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	wg.Wait()
}
