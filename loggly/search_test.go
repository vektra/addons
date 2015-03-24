package loggly

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/cypress"
)

const cAccount = "TEST_LOGGLY_ACCOUNT"
const cUsername = "TEST_LOGGLY_USERNAME"
const cPassword = "TEST_LOGGLY_PASSWORD"

func TestLogglyAPIClientGenerate(t *testing.T) {
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

	ro := &RSIDOptions{}
	eo := &EventsOptions{}
	api, err := NewAPIClient(account, username, password, ro, eo, 100)
	if err != nil {
		panic(err)
	}

	message, err := api.Generate()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, &cypress.Message{}, message)
}
