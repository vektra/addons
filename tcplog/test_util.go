package tcplog

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"testing"

	"github.com/vektra/components/log"
)

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz "
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func NewLogMessage(t *testing.T) *log.Message {
	logMessage := log.Log()
	err := logMessage.Add("message", randString(50))
	err = logMessage.AddString("string_key", "I'm a string!")
	err = logMessage.AddInt("int_key", 12)
	err = logMessage.AddBytes("bytes_key", []byte("I'm bytes!"))
	err = logMessage.AddInterval("interval_key", 2, 1)
	if err != nil {
		t.Errorf("Error adding message: %s", err)
	}
	return logMessage
}

type TcpServer struct {
	Port     int
	Ssl      bool
	Address  chan string
	Messages chan []byte
}

func NewTcpServer() *TcpServer {
	return &TcpServer{
		Ssl:      false,
		Address:  make(chan string, 1),
		Messages: make(chan []byte, 1),
	}
}

func (s *TcpServer) Run(host string) {
	var (
		ln  net.Listener
		err error
	)

	if s.Ssl {
		// need tls cert
		config := tls.Config{InsecureSkipVerify: true}
		ln, err = tls.Listen("tcp", "", &config)
	} else {
		ln, err = net.Listen("tcp", "")
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	s.Address <- fmt.Sprintf("%s:%d", host, ln.Addr().(*net.TCPAddr).Port)

	conn, err := ln.Accept()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	s.handleConnection(conn)
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)

	reqLen, err := conn.Read(buf)

	if err != nil {
		fmt.Println(err)
		return
	}

	s.Messages <- buf[0:reqLen]
}
