package ftp

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"

	"github.com/cnnrznn/goftp/model"
)

func ReceiveFile(ops Option) error {
	conn, err := acceptConn(ops.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	meta, err := receiveMetadata(conn)
	if err != nil {
		return err
	}

	meta.Name = ops.Filename

	if err := receiveFile(conn, meta); err != nil {
		return err
	}

	return nil
}

func acceptConn(addr string) (net.Conn, error) {
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("can't listen on %v: %v", addr, err)
	}
	defer ls.Close()

	conn, err := ls.Accept()
	if err != nil {
		return nil, fmt.Errorf("error accepting request: %v", err)
	}

	return conn, nil
}

func receiveMetadata(conn net.Conn) (*model.Meta, error) {
	metaBs := make([]byte, model.METADATA_SIZE)

	n, err := conn.Read(metaBs)
	if err != nil {
		return nil, err
	}
	if n != model.METADATA_SIZE {
		return nil, fmt.Errorf("received unexpected number of metadata bytes(%v)", n)
	}

	meta, err := model.Deserialize(metaBs)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func receiveFile(conn net.Conn, meta *model.Meta) error {
	buf := make([]byte, 4000000)
	nread := 0
	hasher := sha256.New()

	file, err := os.Create(meta.Name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for nread < meta.Size {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}

		nread += n

		nwritten, err := writer.Write(buf[:n])
		if err != nil {
			return err
		}
		if n != nwritten {
			return fmt.Errorf("did not write all bytes to file")
		}

		nwritten, err = hasher.Write(buf[:n])
		if err != nil {
			return err
		}
		if n != nwritten {
			return fmt.Errorf("did not write all bytes to sha256")
		}
	}

	if !meta.ChecksumEquals(hasher.Sum(nil)) {
		return fmt.Errorf("File checksum does not match metadata!")
	}

	return nil
}
