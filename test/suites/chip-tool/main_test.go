package test

import (
	"fmt"
	"log"
	"matter-snap-testing/test/utils"
	"os"
	"testing"
	"time"
)

const chipToolSnap = "chip-tool"

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestCommon(t *testing.T) {
	go func() {
		utils.Exec(t, fmt.Sprintf("sudo chip-tool pairing onnetwork 110 20202021"))
		time.Sleep(10 * time.Second)
	}()
  
	go func() {
		utils.TestNet(t, chipToolSnap, utils.Net{
			StartSnap:     false,
			TestOpenPorts: []string{utils.ServicePort(chipToolSnap)},
		})

		utils.TestPackaging(t, chipToolSnap, utils.Packaging{
			TestSemanticSnapVersion: true,
		})
	}()

	// This is necessary to prevent the main Goroutine from exiting
	// before the other Goroutines finish executing
	select {}
}

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, chipToolSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, chipToolSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(nil, chipToolSnap)
		}
	}

	if utils.LocalServiceSnap() {
		err = utils.SnapInstallFromFile(nil, utils.LocalServiceSnapPath)
	} else {
		err = utils.SnapInstallFromStore(nil, chipToolSnap, utils.ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	return
}
