package tcplog

import (
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/components/log"
	"github.com/vektra/neko"
)

type TestFormatter struct{}

func (tf *TestFormatter) Format(m *log.Message) ([]byte, error) {
	return []byte(m.KVString()), nil
}

func TestWrite(t *testing.T) {
	n := neko.Start(t)

	var (
		l    *Logger
		line = []byte("This is a log line")
	)

	n.Setup(func() {
		l = NewLogger("", false, &TestFormatter{})
	})

	n.It("adds a log line to the pump", func() {
		l.Write(line)

		select {
		case pumpLine := <-l.Pump:
			assert.Equal(t, line, pumpLine)
			assert.Equal(t, 0, l.PumpDropped)
		default:
			t.Fail()
		}
	})

	n.It("adds an error line to the pump if lines were dropped", func() {
		l.PumpDropped = 1
		l.Write(line)

		select {
		case <-l.Pump:
			expected := "The tcplog pump dropped 1 lines"
			actual := <-l.Pump

			assert.True(t, strings.Index(string(actual), expected) != -1)
			assert.Equal(t, 0, l.PumpDropped)
		default:
			t.Fail()
		}
	})

	n.It("does not add a log line and increments dropped counter if pump is full ", func() {
		l.Pump = make(chan []byte, 0)
		l.Write(line)

		select {
		case <-l.Pump:
			t.Fail()
		default:
			assert.Equal(t, 1, l.PumpDropped)
		}
	})

	n.Meow()
}

func TestDialWhenSslFalse(t *testing.T) {
	s := NewTcpServer()

	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false, &TestFormatter{})

	conn, _ := l.Dial()
	_, ok := conn.(net.Conn)

	assert.True(t, ok, "returns a connection")
}

// func TestDialWhenSslTrue(t *testing.T) {
// 	s := NewTcpServer()
// 	s.Ssl = true

// 	go s.Run()

// 	l := NewLogger(<-s.Address, true, &TestFormatter{})

// 	conn, _ := l.Dial()
// 	_, ok := conn.(net.Conn)

// 	assert.True(t, ok, "returns a connection")
// }

func TestSendLogs(t *testing.T) {
	n := neko.Start(t)

	var (
		s    *TcpServer
		l    *Logger
		line = []byte("This is a log line")
		wg   sync.WaitGroup
	)

	n.Setup(func() {
		s = NewTcpServer()

		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Run("127.0.0.1")
		}()

		l = NewLogger(<-s.Address, false, &TestFormatter{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			l.SendLogs()
		}()
	})

	n.It("sends line from pipe to tcp server", func() {
		l.Pump <- line
		close(l.Pump)

		wg.Wait()

		select {
		case message := <-s.Messages:
			assert.Equal(t, string(line), string(message))
		default:
			t.Fail()
		}
	})

	// n.It("increments dropped counter if tcp write fails", func() {
	// 	l.ConnDropped = 0

	// 	// make write fail ?

	// 	l.Pump <- line
	// 	close(l.Pump)

	// 	wg.Wait()

	// 	assert.Equal(t, 1, l.ConnDropped)
	// })

	n.Meow()
}
