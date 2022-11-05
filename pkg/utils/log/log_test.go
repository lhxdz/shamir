package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	logInfo    = "this is a log info"
	bufferSize = 1024
	tmpLogFile = "/var/log/tmp_shamir_log_test.log"
)

func TestLog(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	SetGlobalLogger(NewLogger(WithLogLever(InfoLevel), WithConsole(), WithLogPath(tmpLogFile)))
	defer func() {
		_ = os.RemoveAll(tmpLogFile)
	}()
	Info(logInfo)

	buffer := make([]byte, bufferSize)
	_, err := r.Read(buffer)
	require.NoError(t, err)
	assert.Contains(t, string(buffer), logInfo)

	fileData, err := os.ReadFile(tmpLogFile)
	require.NoError(t, err)
	assert.Contains(t, string(fileData), logInfo)
}
