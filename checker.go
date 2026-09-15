package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type StateTracker struct {
	State      string
	ChangeTime time.Time
}

type ServiceCurrentState struct {
	mu           sync.Mutex
	ServiceState map[string]StateTracker
}

func (scs *ServiceCurrentState) GetCurrentState(service string) (bool, StateTracker) {

	scs.mu.Lock()
	defer scs.mu.Unlock()
	cs, ok := scs.ServiceState[service]
	return ok, cs
}

func (scs *ServiceCurrentState) UpdateCurrentState(service string, state string) {

	scs.mu.Lock()
	defer scs.mu.Unlock()
	cs, ok := scs.ServiceState[service]

	if !ok || cs.State != state {
		cs.ChangeTime = time.Now()
	}
	cs.State = state
	scs.ServiceState[service] = cs
}

func CheckConnection(host string, port int, timeout time.Duration) (bool, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		return false, err
	}

	conn.Close()
	return true, nil
}

func (scs *ServiceCurrentState) CheckServices(services []ServiceConfig) {
	var wg sync.WaitGroup

	for _, service := range services {
		wg.Go(func() {
			status, err := CheckConnection(service.Host, service.Port, 1*time.Second)
			if err != nil {
				fmt.Printf("Failed with the following error: %v", err)
			}
			if !status {
				scs.UpdateCurrentState(service.Name, "down")

			} else {
				scs.UpdateCurrentState(service.Name, "up")
			}
		})
	}
	wg.Wait()

}
