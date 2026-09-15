package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	cnfg, err := readConfig("./config.json")
	if err != nil {
		fmt.Printf("Reading config failed with: %v", err)
		os.Exit(1)
	}
	scs := &ServiceCurrentState{ServiceState: make(map[string]StateTracker)}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	checkAndPrintServices(cnfg, scs)
	for range ticker.C {
		checkAndPrintServices(cnfg, scs)
	}
}

func checkAndPrintServices(config []ServiceConfig, scs *ServiceCurrentState) {
	scs.CheckServices(config)
	for _, c := range config {
		_, state := scs.GetCurrentState(c.Name)
		fmt.Printf("Current state for %v, is %v for duration %v", c.Name, state.State, time.Since(state.ChangeTime))
	}
}
