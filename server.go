package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Server struct {
	config   *Config
	listener net.Listener
	logger   *log.Logger
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		logger: log.New(os.Stdout, "server : ", log.LstdFlags),
	}
}

func (s *Server) Start() error {
	address := ":" + s.config.Port
	var err error

	if s.config.UseHttps {
		return fmt.Errorf("HTTPS not implemented")
	} else {
		err = s.startHTTP(address)
	}

	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func (s *Server) startHTTP(address string) error {
	s.logger.Println("Starting Http server on port", s.config.Port)
	listener, err := net.Listen("tcp", address)
	s.listener = listener
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	s.logger.Println("HTTP server is ready to accept connections")
	s.acceptConnections()
	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Println("Error accepting connection", err)
			continue
		}
		go s.handleConnection(conn)

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.logger.Println("Closing connection :", conn.RemoteAddr())
		conn.Close()
	}()

	s.logger.Println("Serving new connection :", conn.RemoteAddr())
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
	}
	header := strings.Split(string(buf[:n]), "\r\n")
	if len(header) < 1 {
		s.logger.Println("Invalid request")
	}
	req := header[0]

	if strings.HasPrefix(req, "GET") {
		s.handleGet(conn, req)
	} else if strings.HasPrefix(req, "POST") {
		s.handlePost(conn, req)
	} else {
		s.logger.Println("Unsupported method")
	}
}

func (s *Server) handleGet(conn net.Conn, req string) {
	s.logger.Println("Handling GET request")
}

func (s *Server) handlePost(conn net.Conn, req string) {
	s.logger.Println("Handling POST request")
}
