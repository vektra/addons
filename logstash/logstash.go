package logstash

import (
	"encoding/json"

	"github.com/vektra/addons/lib/tcplog"
	"github.com/vektra/cypress"
	"github.com/vektra/tai64n"
)

const cNewline = "\n"

type LogstashFormatter struct{}

type Message struct {
	Timestamp  *tai64n.TAI64N        `json:"@timestamp,omitempty"`
	Type       *uint32               `json:"type,omitempty"`
	Attributes map[string]*Attribute `json:"attributes,omitempty"`
	SessionId  *string               `json:"session_id,omitempty"`
	Message    string                `json:"message,omitempty"`
}

type Attribute cypress.Attribute

func NewLogger(address string, ssl bool) *tcplog.Logger {
	return tcplog.NewLogger(address, ssl, &LogstashFormatter{})
}

func NewMessage(m *cypress.Message) *Message {
	msg := Message{
		Timestamp:  m.Timestamp,
		Type:       m.Type,
		SessionId:  m.SessionId,
		Attributes: make(map[string]*Attribute),
	}

	for _, a := range m.Attributes {
		if a.GetKey() == cypress.PresetKeys["message"] {
			msg.Message = a.GetSval()
			continue
		}

		attr := Attribute(*a)
		msg.Attributes[a.StringKey()] = &attr
	}

	return &msg
}

func (lf *LogstashFormatter) Format(m *cypress.Message) ([]byte, error) {
	msg := NewMessage(m)

	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	bytes = append(bytes, []byte(cNewline)...)

	return bytes, nil
}
