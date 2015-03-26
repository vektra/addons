package loggly

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vektra/addons/lib/tcplog"
)

const cAccount = "TEST_LOGGLY_ACCOUNT"
const cUsername = "TEST_LOGGLY_USERNAME"
const cPassword = "TEST_LOGGLY_PASSWORD"

func TestLogglyAPIClientGenerate(t *testing.T) {
	token := os.Getenv(cToken)
	if token == "" {
		t.Skipf("%s is not set.", cToken)
	}

	account := os.Getenv(cAccount)
	if account == "" {
		t.Skipf("%s is not set.", cAccount)
	}

	username := os.Getenv(cUsername)
	if username == "" {
		t.Skipf("%s is not set.", cUsername)
	}

	password := os.Getenv(cPassword)
	if password == "" {
		t.Skipf("%s is not set.", cPassword)
	}

	// Send message to loggly

	l := NewLogger(token)

	expected := tcplog.NewMessage(t)
	l.Receive(expected)

	time.Sleep(20 * time.Second)

	// Read back message from loggly

	ro := &RSIDOptions{Size: 1}
	eo := &EventsOptions{}
	api, err := NewAPIClient(account, username, password, ro, eo, 100)
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
