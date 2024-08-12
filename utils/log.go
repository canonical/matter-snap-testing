package utils

import (
	"os"
	"strings"
	"testing"
)

func logFileName(t *testing.T, label string) string {
	fileName := label + ".log"
	if t != nil {
		fileName = strings.ReplaceAll(t.Name(), "/", "-") + "-" + fileName
	}
	logDirectory := "logs"

	err := os.MkdirAll(logDirectory, 0777)
	if err != nil {
		t.Fatalf("Can't create log directory")
	}

	return logDirectory + "/" + fileName
}

func WriteLogFile(t *testing.T, label string, content string) error {
	return os.WriteFile(
		logFileName(t, label),
		[]byte(content),
		0644,
	)
}
