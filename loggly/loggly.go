package loggly

import (
	"bytes"

	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

type LogglyFormatter struct {
	Token string
	PEN   string
}

func NewLogger(address string, ssl bool, token string, pen string) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &LogglyFormatter{token, pen})
}

func (lf *LogglyFormatter) Format(m *cypress.Message) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("<34>1 ") // Need this for loggly to accept. Keep hardcoded?

	buf.WriteString(m.SyslogString(false, false))

	buf.WriteString(" [")
	buf.WriteString(lf.Token)
	buf.WriteString("@")
	buf.WriteString(lf.PEN)
	buf.WriteString("]")

	buf.WriteString("\n")

	return buf.Bytes(), nil
}
