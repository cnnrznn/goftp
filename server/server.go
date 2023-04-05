package server

import (
	"fmt"
	"net"

	"github.com/cnnrznn/goftp/model"
)

type Server struct {
	filename string
	port     int
}

func New(
	port int,
	fn string,
) *Server {
	return &Server{
		port:     port,
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

	meta, err := s.receiveMetadata(conn)
	if err != nil {
		stopChan <- err
		return
	}

	if err := s.receiveFile(conn, meta); err != nil {
		stopChan <- err
		return
	}
}

func (s *Server) acceptConn() (net.Conn, error) {
	ls, err := net.Listen(fmt.Sprintf(":%v", s.port), "tcp")
	if err != nil {
		return nil, fmt.Errorf("can't listen on port %v: %v", s.port, err)
	}
	defer ls.Close()

	conn, err := ls.Accept()
	if err != nil {
		return nil, fmt.Errorf("error accepting request: %v", err)
	}

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

	return meta, nil
}

func (s *Server) receiveFile(conn net.Conn, meta *model.Meta) error {
	// TODO: implement lol
	return nil
}
