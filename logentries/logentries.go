package logentries

import (
	"encoding/json"

	"github.com/vektra/components/lib/tcplog"
	"github.com/vektra/components/log"
)

const cNewline = "\n"

type LogentriesFormatter struct {
	Token string
}

func NewLogger(address string, ssl bool, token string) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &LogentriesFormatter{token})
}

func (lf *LogentriesFormatter) Format(m *log.Message) ([]byte, error) {
	m.Add("token", lf.Token)

	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	bytes = append(bytes, []byte(cNewline)...)

	return bytes, nil
}
