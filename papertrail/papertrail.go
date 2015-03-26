package papertrail

import (
	"encoding/json"

	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

const cNewline = "\n"

type PapertrailFormatter struct{}

func NewLogger(address string, ssl bool) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &PapertrailFormatter{})
}

func (pf *PapertrailFormatter) Format(m *cypress.Message) ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	bytes = append(bytes, []byte(cNewline)...)

	return bytes, nil
}
