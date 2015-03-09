package loggly

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

const cEndpoint = "TEST_LOGGLY_URL"
const cSSL = "TEST_LOGGLY_SSL"
const cToken = "TEST_LOGENTRIES_TOKEN"
const cPEN = "TEST_LOGENTRIES_PEN"

func TestLogglyFormat(t *testing.T) {
	l := NewLogger("", false, "token", "pen")

	message := cypress.Log()
	message.Add("message", "the message")
	message.AddString("string_key", "I'm a string!")
	message.AddInt("int_key", 12)
	message.AddBytes("bytes_key", []byte("I'm bytes!"))
	message.AddInterval("interval_key", 2, 1)

	actual, err := l.Format(message)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	timestamp := message.GetTimestamp().Time().Format(cTimeFormat)

	expected := fmt.Sprintf("%s 0000000 system * message=\"the message\" string_key=\"I'm a string!\" int_key=12 bytes_key=\"I'm bytes!\" interval_key=:2.000000001[token@pen]\n", timestamp)

	assert.Equal(t, expected, string(actual))
}

func TestLogglyRunWithTestServer(t *testing.T) {
	s := tcplog.NewTcpServer()
	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false, "token", "pen")
	go l.Run()

	message := tcplog.NewMessage(t)
	l.Read(message)

	select {
	case m := <-s.Messages:
		expected, err := l.Format(message)
		if err != nil {
			t.Errorf("Error formatting: %s", err)
		}

		assert.Equal(t, string(expected), string(m))

	case <-time.After(5 * time.Second):
		t.Errorf("Test server did not get message in time.")
	}
}

func TestLogglyRunWithLogglyServer(t *testing.T) {
	endpoint := os.Getenv(cEndpoint)
	if endpoint == "" {
		t.Skipf("%s is not set.", cEndpoint)
	}

	ssl := os.Getenv(cSSL)
	if ssl == "" {
		ssl = "true"
	}

	token := os.Getenv(cToken)
	if token == "" {
		t.Skipf("%s is not set.", cToken)
	}

	pen := os.Getenv(cPEN)
	if pen == "" {
		t.Skipf("%s is not set.", cPEN)
	}

	l := NewLogger(endpoint, ssl == "true", token, pen)
	go l.Run()

	message := tcplog.NewMessage(t)
	l.Read(message)

	time.Sleep(10 * time.Second)

	expected, err := l.Format(message)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	t.Logf("Check '%s' got message:\n%s", endpoint, expected)
}
