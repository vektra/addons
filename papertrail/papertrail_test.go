package papertrail

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/components/lib/tcplog"
	"github.com/vektra/components/log"
)

const cEndpoint = "TEST_PAPERTRAIL_URL"
const cSSL = "TEST_PAPERTRAIL_SSL"

func TestPapertrailFormat(t *testing.T) {
	l := NewLogger("", false)

	logMessage := log.Log()
	logMessage.Add("message", "the message")
	logMessage.AddString("string_key", "I'm a string!")
	logMessage.AddInt("int_key", 12)
	logMessage.AddBytes("bytes_key", []byte("I'm bytes!"))
	logMessage.AddInterval("interval_key", 2, 1)

	actual, err := l.Format(logMessage)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	timestamp := logMessage.GetTimestamp().Time().Format(cTimeFormat)

	expected := fmt.Sprintf("%s 0000000 system * message=\"the message\" string_key=\"I'm a string!\" int_key=12 bytes_key=\"I'm bytes!\" interval_key=:2.000000001\n", timestamp)

	assert.Equal(t, expected, string(actual))
}

func TestPapertrailRunWithTestServer(t *testing.T) {
	if !log.Available() {
		t.Skip("Log is not availble.")
	}

	s := tcplog.NewTcpServer()
	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false)
	go l.WatchLogs()
	go l.SendLogs()
	defer l.Cleanup()

	logMessage := tcplog.NewLogMessage(t)
	logMessage.Inject()

	time.Sleep(1 * time.Second)

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

func TestPapertrailRunWithPapertrailServer(t *testing.T) {
	if !log.Available() {
		t.Skip("Log is not available.")
	}

	endpoint := os.Getenv(cEndpoint)
	if endpoint == "" {
		t.Skipf("%s is not set.", cEndpoint)
	}

	ssl := os.Getenv(cSSL)
	if ssl == "" {
		ssl = "true"
	}

	l := NewLogger(endpoint, ssl == "true")
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
