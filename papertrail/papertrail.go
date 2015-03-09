package papertrail

import (
	"bytes"

	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

type PapertrailFormatter struct{}

func NewLogger(address string, ssl bool) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &PapertrailFormatter{})
}

func (pf *PapertrailFormatter) Format(m *cypress.Message) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(m.SyslogString(false, false))

	buf.WriteString("\n")

	return buf.Bytes(), nil
}
