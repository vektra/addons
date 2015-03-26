package logentries

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vektra/addons/lib/tcplog"
)

const cAccountKey = "TEST_LOGENTRIES_ACCOUNT_KEY"
const cHost = "TEST_LOGENTRIES_HOST"
const cLog = "TEST_LOGENTRIES_LOG"

func TestLogentriesAPIClientGenerate(t *testing.T) {
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

	key := os.Getenv(cAccountKey)
	if key == "" {
		t.Skipf("%s is not set.", cAccountKey)
	}

	host := os.Getenv(cHost)
	if host == "" {
		t.Skipf("%s is not set.", cHost)
	}

	log := os.Getenv(cLog)
	if log == "" {
		t.Skipf("%s is not set.", cLog)
	}

	// Send message to logentries

	l := NewLogger(endpoint, ssl == "true", token)
	go l.Run()

	expected := tcplog.NewMessage(t)
	l.Read(expected)

	time.Sleep(10 * time.Second)

	// Read back message from logentries

	timestamp, _ := json.Marshal(expected.Timestamp)

	options := &EventsOptions{Filter: fmt.Sprintf("/%s", timestamp)}
	api, err := NewAPIClient(key, host, log, options, 100)
	require.NoError(t, err)

	actual, err := api.Generate()
	require.NoError(t, err)

	// Make sure its the same message

	require.Equal(t, expected.GetTimestamp().Time(), actual.GetTimestamp().Time())
	require.Equal(t, expected.GetVersion(), actual.GetVersion())
	require.Equal(t, expected.GetSessionId(), actual.GetSessionId())
	require.Equal(t, expected.GetTags(), actual.GetTags())

	expectedMessage, _ := expected.GetString("message")
	actualMessage, _ := actual.GetString("message")

	require.Equal(t, expectedMessage, actualMessage)
}
