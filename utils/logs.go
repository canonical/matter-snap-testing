package utils

import (
	"os"
	"strings"
	"testing"
)

func GetLogFileName(t *testing.T, label string) string {
	fileName := strings.ReplaceAll(t.Name(), "/", "-") + "-" + label + ".log"
	logDirectory := "logs"

	err := os.MkdirAll(logDirectory, 0777)
	if err != nil {
		t.Fatalf("can't create log directory")
	}

	return logDirectory + "/" + fileName
}

func WriteLogFile(t *testing.T, label string, b []byte) error {
	return os.WriteFile(
		GetLogFileName(t, label),
		b,
		0777,
	)
}
