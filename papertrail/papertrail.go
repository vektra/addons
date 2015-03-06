package papertrail

import (
	"bytes"
	"time"

	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/components/lib/request"
	"github.com/vektra/cypress"
)

const cTimeFormat = time.RFC3339Nano
const cNewline = "\n"

type PapertrailFormatter struct{}

func NewLogger(address string, ssl bool) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &PapertrailFormatter{})
}

func (pf *PapertrailFormatter) Format(m *cypress.Message) ([]byte, error) {
	var buf bytes.Buffer

	time := m.GetTimestamp().Time().Format(cTimeFormat)

	buf.WriteString(time)
	buf.WriteString(" ")

	if s := m.GetSessionId(); len(s) > 0 {
		buf.WriteString(request.Id(s).String()[0:7])
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
				buf.WriteString(cNewline)

				return buf.Bytes(), nil
			}
		}
	}

	buf.WriteString(" system ")
	buf.WriteString("*")
	buf.WriteString(m.KVPairs())
	buf.WriteString(cNewline)

	return buf.Bytes(), nil
}
