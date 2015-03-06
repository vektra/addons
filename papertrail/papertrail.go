package papertrail

import (
	"time"

	"github.com/vektra/components/lib/tcplog"
	"github.com/vektra/components/log"
)

const cTimeFormat = time.RFC3339Nano
const cNewline = "\n"

type PapertrailFormatter struct{}

func NewLogger(address string, ssl bool) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &PapertrailFormatter{})
}

func (pf *PapertrailFormatter) Format(m *log.Message) ([]byte, error) {
	return []byte(m.SyslogString(false, false) + cNewline), nil
}
