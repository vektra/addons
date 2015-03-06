package logentries

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/components/lib/tcplog"
	"github.com/vektra/components/log"
)

const cEndpoint = "TEST_LOGENTRIES_URL"
const cSSL = "TEST_LOGENTRIES_SSL"
const cToken = "TEST_LOGENTRIES_TOKEN"

func TestLogentriesFormat(t *testing.T) {
	l := NewLogger("", false, "token")

	logMessage := log.Log()
	logMessage.AddString("string_key", "I'm a string!")
	logMessage.AddInt("int_key", 12)
	logMessage.AddBytes("bytes_key", []byte("I'm bytes!"))
	logMessage.AddInterval("interval_key", 2, 1)

	actual, err := l.Format(logMessage)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	timestamp, err := json.Marshal(logMessage.Timestamp)
	if err != nil {
		t.Errorf("Error marshalling timestamp to JSON: %s", err)
	}

	expected := fmt.Sprintf("{\"timestamp\":%s,\"type\":0,\"attributes\":[{\"string_key\":\"I'm a string!\"},{\"int_key\":12},{\"bytes_key\":\"SSdtIGJ5dGVzIQ==\",\"_bytes\":\"\"},{\"interval_key\":{\"seconds\":2,\"nanoseconds\":1}},{\"token\":\"token\"}]}\n", timestamp)

	assert.Equal(t, expected, string(actual))
}

func TestLogentriesRunWithTestServer(t *testing.T) {
	if !log.Available() {
		t.Skip("Log is not availble.")
	}

	s := tcplog.NewTcpServer()
	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false, "token")
	go l.WatchLogs()
	go l.SendLogs()
	defer l.Cleanup()

	logMessage := tcplog.NewLogMessage(t)
	logMessage.Inject()

	select {
	case message := <-s.Messages:
		expected, err := l.Format(logMessage)
		if err != nil {
			t.Errorf("Error formatting: %s", err)
		}

		assert.Equal(t, string(expected), string(message))

	case <-time.After(5 * time.Second):
		t.Errorf("Test server did not get message in time.")
	}
}

func TestLogentriesRunWithLogentriesServer(t *testing.T) {
	if !log.Available() {
		t.Skip("log is not available.")
	}

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
	go l.WatchLogs()
	go l.SendLogs()
	defer l.Cleanup()

	logMessage := tcplog.NewLogMessage(t)
	logMessage.Inject()

	time.Sleep(10 * time.Second)

	expected, err := l.Format(logMessage)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	t.Logf("Check '%s' got message:\n%s", endpoint, expected)
}
