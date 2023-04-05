package client

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/cnnrznn/goftp/model"
)

type Client struct {
	filename string
	toURI    string
}

func New(
	fn string,
	toURI string,
) *Client {
	return &Client{
		filename: fn,
		toURI:    toURI,
	}
}

func (c *Client) Run(stopChan chan error) {
	meta, err := c.createMetadata()
	if err != nil {
		stopChan <- err
		return
	}

	if err := c.sendFile(meta); err != nil {
		stopChan <- err
		return
	}
}

func (c *Client) createMetadata() (*model.Meta, error) {
	meta := &model.Meta{}

	file, err := os.Open(c.filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file %v", c.filename)
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

func (c *Client) sendFile(meta *model.Meta) error {
	return nil
}

func calculateChecksum(file *os.File, size int) (string, error) {
	hasher := sha256.New()
	buf := make([]byte, 8192)
	reader := bufio.NewReader(file)
	nread := 0

	for nread < size {
		n, err := reader.Read(buf)
		if err != nil {
			return "", err
		}

		nread += n
		nwritten, err := hasher.Write(buf[:n])
		if err != nil {
			return "", err
		}
		if n != nwritten {
			return "", fmt.Errorf("failed to write all bytes to hasher")
		}
	}

	return string(hasher.Sum(nil)), nil
}
