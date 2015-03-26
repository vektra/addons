package papertrail

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vektra/addons/lib/tcplog"
)

const cToken = "TEST_PAPERTRAIL_TOKEN"

func TestPapertrailAPIGenerate(t *testing.T) {
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

	// Send message to papertrail

	l := NewLogger(endpoint, ssl == "true")
	go l.Run()

	expected := tcplog.NewMessage(t)
	l.Receive(expected)

	time.Sleep(10 * time.Second)

	// Read back message from papertrail

	options := &EventsOptions{}
	api := NewAPIClient(token, options, 100)

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
