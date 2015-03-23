package logentries

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/cypress"
)

const cAccountKey = "TEST_LOGENTRIES_ACCOUNT_KEY"
const cHost = "TEST_LOGENTRIES_HOST"
const cLog = "TEST_LOGENTRIES_LOG"

func TestLogentriesAPIClientGenerate(t *testing.T) {
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

	options := &EventsOptions{}
	api, err := NewAPIClient(key, host, log, options, 100)
	if err != nil {
		panic(err)
	}

	message, err := api.Generate()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, &cypress.Message{}, message)
}
