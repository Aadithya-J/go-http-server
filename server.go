package main

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Server struct {
	config   *Config
	listener net.Listener
	logger   *log.Logger
	wg       sync.WaitGroup
	quit     chan struct{}
	limiter  *Limiter
}

func NewServer(config *Config) *Server {
	return &Server{
		config:  config,
		logger:  log.New(os.Stdout, "server : ", log.LstdFlags),
		quit:    make(chan struct{}),
		limiter: NewLimiter(),
	}
}

func (s *Server) Start() error {
	address := ":" + s.config.Port
	var err error

	if s.config.UseHttps {
		err = s.startHTTPS(address)
	} else {
		err = s.startHTTP(address)
	}

	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

func (s *Server) startHTTP(address string) error {
	s.logger.Println("Starting Http server on port", s.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	s.listener = listener
	s.logger.Println("HTTP server is ready to accept connections")
	return nil
}

func (s *Server) startHTTPS(address string) error {
	s.logger.Println("Starting Https server on port", s.config.Port)
	cert, err := tls.LoadX509KeyPair(s.config.CertFile, s.config.KeyFile)
	if err != nil {
		return fmt.Errorf("Error loading certificate: %v", err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return fmt.Errorf("Error starting HTTPS server: %v", err)
	}
	s.listener = listener
	s.logger.Println("HTTPS server is ready to accept connections")
	return nil
}

func (s *Server) Stop() error {
	close(s.quit)

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.logger.Println("Error closing listener:", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("server shutdown timed out")
	case <-done:
		return nil
	}
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					s.logger.Println("Error accepting connection", err)
					continue
				}
			}
			clientIP := extractClientIP(conn.RemoteAddr().String())

			if !s.limiter.Allow(clientIP, s) {
				s.logger.Printf("Rate limit exceeded for client %s\n", clientIP)
				conn.Close()
				continue
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.logger.Println("Closing connection:", conn.RemoteAddr())
		conn.Close()
		s.wg.Done()
	}()

	s.logger.Println("Serving new connection:", conn.RemoteAddr())
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		s.logger.Println("Error reading:", err)
		return
	}

	header := strings.Split(string(buf[:n]), "\r\n")
	if len(header) < 1 {
		s.logger.Println("Invalid request")
		return
	}
	req := header[0]

	if strings.HasPrefix(req, "GET") {
		s.handleGet(conn, req)
	} else {
		s.logger.Println("Unsupported method")
	}
}

func (s *Server) handleGet(conn net.Conn, req string) {

	path := req[4 : strings.Index(req, "HTTP")-1]
	s.logger.Println("Request path:", path)
	if strings.Contains(path, ".") {
		s.serveFiles(conn, path)
	} else if path == "/home" {
		s.serveFiles(conn, "index.html")
	} else {
		write200(conn, "text/plain", "Hello, World!")
	}
}

func (s *Server) serveFiles(conn net.Conn, path string) {
	fileDir := filepath.Join(s.config.Directory, path)
	s.logger.Println("Serving file:", fileDir)
	fileData, err := readFile(fileDir)
	if err != nil {
		s.logger.Println("Error reading file:", err)
		write400(conn)
		return
	}
	ext := filepath.Ext(path)
	contentType := getContentType(ext)

	if strings.Contains(conn.RemoteAddr().String(), "gzip") {
		write200(conn, contentType, "")
		gw := gzip.NewWriter(conn)
		defer gw.Close()

		if _, err := gw.Write(fileData); err != nil {
			log.Fatal("Gzip error:", err)
		}
	} else {
		write200(conn, contentType, string(fileData))
	}
}

func readFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getContentType(ext string) string {
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	default:
		return "text/plain"
	}
}

func write200(conn net.Conn, textType string, body string) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", textType, len(body), body)
	conn.Write([]byte(response))
}

func write400(conn net.Conn) {
	response := "HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"
	conn.Write([]byte(response))
}

func (s *Server) static(dir string) {
	s.config.Directory = dir
}
