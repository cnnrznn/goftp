package server

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"

	"github.com/cnnrznn/goftp/model"
)

type Server struct {
	filename string
	addr     string
}

func New(
	addr string,
	fn string,
) *Server {
	return &Server{
		addr:     addr,
		filename: fn,
	}
}

func (s *Server) Run(stopChan chan error) {
	conn, err := s.acceptConn()
	if err != nil {
		stopChan <- err
		return
	}
	defer conn.Close()

	fmt.Println("receiving metadata")

	meta, err := s.receiveMetadata(conn)
	if err != nil {
		stopChan <- err
		return
	}

	fmt.Println("receving file")

	if err := s.receiveFile(conn, meta); err != nil {
		stopChan <- err
		return
	}
}

func (s *Server) acceptConn() (net.Conn, error) {
	fmt.Printf("Listening on %v\n", s.addr)

	ls, err := net.Listen("tcp", s.addr)
	if err != nil {
		return nil, fmt.Errorf("can't listen on %v: %v", s.addr, err)
	}
	defer ls.Close()

	conn, err := ls.Accept()
	if err != nil {
		return nil, fmt.Errorf("error accepting request: %v", err)
	}

	fmt.Println("Accepted connection")
	return conn, nil
}

func (s *Server) receiveMetadata(conn net.Conn) (*model.Meta, error) {
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

	fmt.Println("received metadata")

	return meta, nil
}

func (s *Server) receiveFile(conn net.Conn, meta *model.Meta) error {
	buf := make([]byte, 4000000)
	nread := 0
	hasher := sha256.New()

	file, err := os.Open(meta.Name)
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

	if meta.Checksum != string(hasher.Sum(nil)) {
		return fmt.Errorf("File checksum does not match metadata!")
	}

	return nil
}
