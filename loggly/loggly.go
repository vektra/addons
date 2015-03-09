package loggly

import (
	"bytes"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

const cTimeFormat = time.RFC3339Nano
const cNewline = "\n"

type LogglyFormatter struct {
	Token string
	PEN   string
}

func NewLogger(address string, ssl bool, token string, pen string) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &LogglyFormatter{token, pen})
}

func (lf *LogglyFormatter) Format(m *cypress.Message) ([]byte, error) {
	var buf bytes.Buffer

	time := m.GetTimestamp().Time().Format(cTimeFormat)

	buf.WriteString(time)
	buf.WriteString(" ")

	if s := m.GetSessionId(); len(s) > 0 {
		buf.WriteString(uuid.UUID(s).String()[0:7])
	} else {
		buf.WriteString("0000000")
	}

	// Special case the logs that come out of the volts to make them easier to read
	if m.GetType() == 0 {
		if volt, ok := m.GetString("volt"); ok {
			if log, ok := m.GetString("log"); ok {

				buf.WriteString(" ")
				buf.WriteString(volt)
				buf.WriteString(" ")
				buf.WriteString(log)

				buf.WriteString("[")
				buf.WriteString(lf.Token)
				buf.WriteString("@")
				buf.WriteString(lf.PEN)
				buf.WriteString("]")

				buf.WriteString(cNewline)

				return buf.Bytes(), nil
			}
		}
	}

	buf.WriteString(" system ")
	buf.WriteString("*")
	buf.WriteString(m.KVPairs())

	buf.WriteString("[")
	buf.WriteString(lf.Token)
	buf.WriteString("@")
	buf.WriteString(lf.PEN)
	buf.WriteString("]")

	buf.WriteString(cNewline)

	return buf.Bytes(), nil
}
