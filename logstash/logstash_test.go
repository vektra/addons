package logstash

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/components/lib/tcplog"
	"github.com/vektra/components/log"
)

const cLogstash = "/usr/local/logstash-1.4.2/bin/logstash"

func TestLogstashFormat(t *testing.T) {
	l := NewLogger("", false)

	logMessage := log.Log()
	logMessage.Add("message", "the message")
	logMessage.AddString("string_key", "I'm a string!")
	logMessage.AddInt("int_key", 12)
	logMessage.AddBytes("bytes_key", []byte("I'm bytes!"))
	logMessage.AddInterval("interval_key", 2, 1)

	actual, err := l.Format(logMessage)
	if err != nil {
		t.Errorf("Error formatting: %s", err)
	}

	timestamp, err := json.Marshal(logMessage.Timestamp)
	if err != nil {
		t.Errorf("Error marshalling timestamp to JSON: %s", err)
	}

	expected := fmt.Sprintf("{\"@timestamp\":%s,\"type\":0,\"attributes\":{\"bytes_key\":{\"value\":\"SSdtIGJ5dGVzIQ==\",\"_bytes\":\"\"},\"int_key\":12,\"interval_key\":{\"seconds\":2,\"nanoseconds\":1},\"string_key\":\"I'm a string!\"},\"message\":\"the message\"}\n", timestamp)

	assert.Equal(t, expected, string(actual))
}

func TestLogstashRunWithTestServer(t *testing.T) {
	if !log.Available() {
		t.Skip("Log is not availble.")
	}

	s := tcplog.NewTcpServer()
	go s.Run("127.0.0.1")

	l := NewLogger(<-s.Address, false)
	go l.WatchLogs()
	go l.SendLogs()
	defer l.Cleanup()

	logMessage := tcplog.NewLogMessage(t)
	logMessage.Inject()

	select {
	case message := <-s.Messages:
		expected, err := l.Format(logMessage)
		if err != nil {
			t.Errorf("Error formatting: %s", err)
		}

		assert.Equal(t, string(expected), string(message))

	case <-time.After(5 * time.Second):
		t.Errorf("Test server did not get message in time.")
	}
}

func TestLogstashRunWithLogstashServer(t *testing.T) {
	if !log.Available() {
		t.Skip("Log is not available.")
	}

	// Check for logstash
	if _, err := os.Stat(cLogstash); err != nil {
		t.Skip("Logstash is not available.")
	}

	// Find free port
	ln, err := net.Listen("tcp", "")
	if err != nil {
		t.Fatal(err)
	}
	port := fmt.Sprintf("%d", ln.Addr().(*net.TCPAddr).Port)
	ln.Close()

	// Start logstash on found port
	config := fmt.Sprintf("input { tcp { port => %s codec => json_lines {} } } output { stdout {} }", port)
	cmd := exec.Command("bin/logstash", "-e", config)
	cmd.Env = []string{"PATH=/usr/local/bin:/usr/bin:/usr/sbin:/sbin:/bin"}
	cmd.Dir = "/usr/local/logstash-1.4.2"
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	defer cmd.Process.Kill()

	time.Sleep(1 * time.Second)

	// Send logs to logstash on found port
	l := NewLogger("0.0.0.0:"+port, false)
	go l.WatchLogs()
	go l.SendLogs()
	defer l.Cleanup()

	logMessage := tcplog.NewLogMessage(t)
	logMessage.Inject()

	time.Sleep(1 * time.Second)

	msg := NewMessage(logMessage)
	expected := msg.Message

	r := bufio.NewReader(stdout)
	out, _, err := r.ReadLine() // throw away first line
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("first line: %s\n", string(out))

	out, _, err = r.ReadLine()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, strings.Index(string(out), string(expected)) != -1,
		fmt.Sprintf("Expected: %s Got: %s", expected, string(out)))
}
