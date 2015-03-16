package logentries

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cAccountKey = "TEST_LOGENTRIES_ACCOUNT_KEY"
const cHost = "TEST_LOGENTRIES_HOST"
const cLog = "TEST_LOGENTRIES_LOG"

func TestLogentriesAPIClientSearch(t *testing.T) {
	accountKey := os.Getenv(cAccountKey)
	if accountKey == "" {
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

	api := &APIClient{accountKey, host, log}
	options := &EventsOptions{}

	actual, err := api.Search(options)

	assert.NoError(t, err)
	assert.Equal(t, "", actual)
}
