package loggly

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cAccount = "TEST_LOGGLY_ACCOUNT"
const cUsername = "TEST_LOGGLY_USERNAME"
const cPassword = "TEST_LOGGLY_PASSWORD"

func TestLogglyAPIClientSearch(t *testing.T) {
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

	api := &APIClient{account, username, password}
	ro := &RSIDOptions{}
	eo := &EventsOptions{}

	actual, err := api.Search(ro, eo)

	assert.NoError(t, err)
	assert.Equal(t, "", actual)
}
