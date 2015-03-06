package logstash

import (
	"encoding/json"

	"github.com/vektra/components/lib/tcplog"
	"github.com/vektra/components/log"
)

const cNewline = "\n"

type LogstashFormatter struct{}

type Message struct {
	Timestamp  *log.TAI64N           `json:"@timestamp,omitempty"`
	Type       *uint32               `json:"type,omitempty"`
	Attributes map[string]*Attribute `json:"attributes,omitempty"`
	SessionId  []byte                `json:"session_id,omitempty"`
	Message    string                `json:"message,omitempty"`
}

type Attribute log.Attribute

func NewLogger(address string, ssl bool) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &LogstashFormatter{})
}

func NewMessage(m *log.Message) *Message {
	msg := Message{
		Timestamp:  m.Timestamp,
		Type:       m.Type,
		SessionId:  m.SessionId,
		Attributes: make(map[string]*Attribute),
	}

	for _, a := range m.Attributes {
		if a.GetKey() == log.PresetKeys["message"] {
			msg.Message = a.GetSval()
			continue
		}

		attr := Attribute(*a)
		msg.Attributes[a.StringKey()] = &attr
	}

	return &msg
}

func (lf *LogstashFormatter) Format(m *log.Message) ([]byte, error) {
	msg := NewMessage(m)

	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	bytes = append(bytes, []byte(cNewline)...)

	return bytes, nil
}
