package logentries

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

const cEndpoint = "TEST_LOGENTRIES_URL"
const cSSL = "TEST_LOGENTRIES_SSL"
const cToken = "TEST_LOGENTRIES_TOKEN"

func TestLogentriesFormat(t *testing.T) {
	l := NewLogger("", false, "token")

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

	timestamp, err := json.Marshal(message.Timestamp)
	if err != nil {
		t.Errorf("Error marshalling timestamp to JSON: %s", err)
	}

	expected := fmt.Sprintf("{\"@timestamp\":%s,\"@type\":\"log\",\"@version\":\"1\",\"bytes_key\":{\"bytes\":\"SSdtIGJ5dGVzIQ==\"},\"int_key\":12,\"interval_key\":{\"nanoseconds\":1,\"seconds\":2},\"message\":\"the message\",\"string_key\":\"I'm a string!\",\"token\":\"token\"}\n", timestamp)

	assert.Equal(t, expected, string(actual))
}

func TestLogentriesRunWithTestServer(t *testing.T) {
	s := tcplog.NewTcpServer()
	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false, "token")
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

func TestLogentriesRunWithLogentriesServer(t *testing.T) {
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

	l := NewLogger(endpoint, ssl == "true", token)
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
