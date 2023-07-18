package ftp

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/cnnrznn/goftp/model"
)

func SendFile(ops Option) error {
	fn := ops.Filename
	destination := ops.Addr

	meta, err := model.GetMetadata(fn)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendMetadata(meta, conn); err != nil {
		return err
	}

	if err := sendFile(meta, conn); err != nil {
		return err
	}

	return nil
}

func sendMetadata(meta *model.Meta, conn net.Conn) error {
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

func sendFile(meta *model.Meta, conn net.Conn) error {
	file, err := os.Open(meta.Name)
	if err != nil {
		return err
	}
	defer file.Close()

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
