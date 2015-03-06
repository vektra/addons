package logentries

import (
	"encoding/json"

	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
)

const cNewline = "\n"

type LogentriesFormatter struct {
	Token string
}

func NewLogger(address string, ssl bool, token string) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &LogentriesFormatter{token})
}

func (lf *LogentriesFormatter) Format(m *cypress.Message) ([]byte, error) {
	if _, ok := m.Get("token"); !ok {
		m.Add("token", lf.Token)
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	bytes = append(bytes, []byte(cNewline)...)

	return bytes, nil
}
