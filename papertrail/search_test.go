package papertrail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cToken = "TEST_PAPERTRAIL_TOKEN"

func TestPapertrailAPIClientSearch(t *testing.T) {
	token := os.Getenv(cToken)
	if token == "" {
		t.Skipf("%s is not set.", cToken)
	}

	api := &APIClient{token}
	options := &EventsOptions{}

	actual, err := api.Search(options)

	assert.NoError(t, err)
	assert.Equal(t, "", actual)
}
