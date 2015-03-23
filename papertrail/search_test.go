package papertrail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/cypress"
)

const cToken = "TEST_PAPERTRAIL_TOKEN"

func TestPapertrailAPIGenerate(t *testing.T) {
	token := os.Getenv(cToken)
	if token == "" {
		t.Skipf("%s is not set.", cToken)
	}

	options := &EventsOptions{}
	api := NewAPIClient(token, options, 100)

	message, err := api.Generate()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, &cypress.Message{}, message)
}
