package client

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
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

	conn, err := net.Dial("tcp", c.toURI)
	if err != nil {
		stopChan <- err
		return
	}
	defer conn.Close()

	if err := c.sendMetadata(meta, conn); err != nil {
		stopChan <- err
		return
	}

	if err := c.sendFile(meta, conn); err != nil {
		stopChan <- err
		return
	}

	stopChan <- nil
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

func calculateChecksum(file *os.File, size int) ([]byte, error) {
	hasher := sha256.New()
	buf := make([]byte, 8192)
	reader := bufio.NewReader(file)
	nread := 0

	for nread < size {
		n, err := reader.Read(buf)
		if err != nil {
			return []byte{}, err
		}

		nread += n
		nwritten, err := hasher.Write(buf[:n])
		if err != nil {
			return []byte{}, err
		}
		if n != nwritten {
			return []byte{}, fmt.Errorf("failed to write all bytes to hasher")
		}
	}

	return hasher.Sum(nil), nil
}

func (c *Client) sendMetadata(meta *model.Meta, conn net.Conn) error {
	bs, err := model.Serialize(meta)
	if err != nil {
		return err
	}

	sent, err := conn.Write(bs)
	if err != nil {
		return err
	}
	if sent != len(bs) {
		return fmt.Errorf("did not send correct metadata bytes")
	}

	return nil
}

func (c *Client) sendFile(meta *model.Meta, conn net.Conn) error {
	file, err := os.Open(meta.Name)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)

	buf := make([]byte, 64000)

	for nread := 0; nread < meta.Size; {
		n, err := reader.Read(buf)
		if err != nil {
			return err
		}

		nread += n

		nwritten, err := conn.Write(buf[:n])
		if err != nil {
			return err
		}
		if nwritten != n {
			return fmt.Errorf("not all bytes written to server")
		}
	}

	return nil
}
