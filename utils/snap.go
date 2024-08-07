package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// func SnapInstall(t *testing.T, name string) {
// 	if strings.HasSuffix(name, ".snap") {
// 		SnapInstallFromFile(nil, name)
// 	} else {
// 		SnapInstallFromStore(nil, name, ServiceChannel)
// 	}
// }

func SnapInstallFromStore(t *testing.T, name, channel string) error {

	option := "--channel"
	// install by revision if channel is a number
	if _, err := strconv.Atoi(channel); err == nil {
		option = "--revision"
	}

	_, stderr, err := ExecVerbose(t, fmt.Sprintf(
		"sudo snap install %s %s=%s",
		name,
		option,
		channel,
	))

	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapInstallFromFile(t *testing.T, path string) error {
	_, stderr, err := ExecVerbose(t, fmt.Sprintf(
		"sudo snap install --dangerous %s",
		path,
	))
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapInstalled(t *testing.T, name string) bool {
	out, _, _ := ExecVerbose(t, fmt.Sprintf(
		"snap list %s || true",
		name,
	))
	return strings.Contains(out, name)
}

func SnapRemove(t *testing.T, names ...string) {
	for _, name := range names {
		ExecVerbose(t, fmt.Sprintf(
			"sudo snap remove --purge %s",
			name,
		))
	}
}

func SnapBuild(t *testing.T, workDir string) error {
	_, stderr, err := ExecVerbose(t, fmt.Sprintf(
		"cd %s && snapcraft",
		workDir,
	))
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapConnect(t *testing.T, plug, slot string) error {
	_, stderr, err := ExecVerbose(t, fmt.Sprintf(
		"sudo snap connect %s %s",
		plug, slot,
	))
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapConnectSecretstoreToken(t *testing.T, snap string) error {
	return SnapConnect(t,
		"edgexfoundry:edgex-secretstore-token",
		snap+":edgex-secretstore-token")
}

func SnapDisconnect(t *testing.T, plug, slot string) {
	ExecVerbose(t, fmt.Sprintf(
		"sudo snap disconnect %s %s",
		plug, slot,
	))
}

func SnapVersion(t *testing.T, name string) string {
	out, _, _ := ExecVerbose(t, fmt.Sprintf(
		"snap info %s | grep installed | awk '{print $2}'",
		name,
	))
	return strings.TrimSpace(out)
}

func SnapRevision(t *testing.T, name string) string {
	out, _, _ := ExecVerbose(t, fmt.Sprintf(
		"snap list %s | awk 'NR==2 {print $3}'",
		name,
	))
	return strings.TrimSpace(out)
}

func snapJournalCommand(start time.Time, name string) string {
	// The command should not return error even if nothing is grepped, hence the "|| true"
	return fmt.Sprintf("sudo journalctl --since \"%s\" --no-pager | grep \"%s\"|| true",
		start.Format("2006-01-02 15:04:05"),
		name)
}

func SnapDumpLogs(t *testing.T, start time.Time, snapName string) {
	logFileName := logFileName(t, snapName)

	ExecVerbose(t, fmt.Sprintf("(%s) > %s",
		snapJournalCommand(start, snapName),
		logFileName))

	wd, _ := os.Getwd()
	fmt.Printf("Wrote snap logs to %s/%s\n", wd, logFileName)
}

func SnapLogs(t *testing.T, start time.Time, name string) string {
	logs, _, _ := Exec(t, snapJournalCommand(start, name))
	return logs
}

func SnapSet(t *testing.T, name, key, value string) {
	ExecVerbose(t, fmt.Sprintf(
		"sudo snap set %s %s='%s'",
		name,
		key,
		value,
	))
}

func SnapUnset(t *testing.T, name string, keys ...string) {
	ExecVerbose(t, fmt.Sprintf(
		"sudo snap unset %s %s",
		name,
		strings.Join(keys, " "),
	))
}

func SnapStart(t *testing.T, names ...string) {
	for _, name := range names {
		ExecVerbose(t, fmt.Sprintf(
			"sudo snap start --enable %s",
			name,
		))
	}
}

func SnapStop(t *testing.T, names ...string) {
	for _, name := range names {
		ExecVerbose(t, fmt.Sprintf(
			"sudo snap stop --disable %s",
			name,
		))
	}
}

func SnapRestart(t *testing.T, names ...string) {
	for _, name := range names {
		ExecVerbose(t, fmt.Sprintf(
			"sudo snap restart %s",
			name,
		))
	}
	// Add delay after restart to avoid reaching systemd's restart limits
	// See https://www.freedesktop.org/software/systemd/man/systemd-system.conf.html#DefaultStartLimitIntervalSec=
	time.Sleep(1 * time.Second)
}

func SnapRefresh(t *testing.T, name, channel string) {
	ExecVerbose(t, fmt.Sprintf(
		"sudo snap refresh %s --channel=%s --amend",
		name,
		channel,
	))
}

func SnapServicesEnabled(t *testing.T, name string) bool {
	out, _, _ := ExecVerbose(t, fmt.Sprintf(
		"snap services %s | awk 'FNR == 2 {print $2}'",
		name,
	))
	return strings.TrimSpace(out) == "enabled"
}

func SnapServicesActive(t *testing.T, name string) bool {
	out, _, _ := ExecVerbose(t, fmt.Sprintf(
		"snap services %s | awk 'FNR == 2 {print $3}'",
		name,
	))
	return strings.TrimSpace(out) == "active"
}
