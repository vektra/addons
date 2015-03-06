package papertrail

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

const cEndpoint = "TEST_PAPERTRAIL_URL"
const cSSL = "TEST_PAPERTRAIL_SSL"

func TestPapertrailFormat(t *testing.T) {
	l := NewLogger("", false)

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

	expected := fmt.Sprintf("%s 0000000 system * message=\"the message\" string_key=\"I'm a string!\" int_key=12 bytes_key=\"I'm bytes!\" interval_key=:2.000000001\n", timestamp)

	assert.Equal(t, expected, string(actual))
}

func TestPapertrailRunWithTestServer(t *testing.T) {
	s := tcplog.NewTcpServer()
	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false)
	go l.SendLogs()
	defer l.Cleanup()

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

func TestPapertrailRunWithPapertrailServer(t *testing.T) {
	endpoint := os.Getenv(cEndpoint)
	if endpoint == "" {
		t.Skipf("%s is not set.", cEndpoint)
	}

	ssl := os.Getenv(cSSL)
	if ssl == "" {
		ssl = "true"
	}

	l := NewLogger(endpoint, ssl == "true")
	go l.SendLogs()
	defer l.Cleanup()

	message := tcplog.NewMessage(t)
	l.Read(message)

	time.Sleep(10 * time.Second)

	expected, err := l.Format(message)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	t.Logf("Check '%s' got message:\n%s", endpoint, expected)
}
