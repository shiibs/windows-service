package main

import (
	"log"
	"os/exec"
	"testing"
	"time"
)

// TestInstallService verifies that the service can be installed correctly.
func TestInstallService(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "install")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to install service: %v", err)
	}

	// Verify installation
	cmd = exec.Command("sc", "query", serviceName)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Service installation verification failed: %v", err)
	}
	log.Println("Service installed.")
}

// TestStartService verifies that the service can be started correctly.
func TestStartService(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "start")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}

	// Wait for a bit to allow the service to start
	time.Sleep(10 * time.Second)

	// Verify service is running
	cmd = exec.Command("sc", "query", serviceName)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Service start verification failed: %v", err)
	}
	log.Println("Service started.")
}

// TestStopService verifies that the service can be stopped correctly.
func TestStopService(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "stop")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to stop service: %v", err)
	}

	// Verify service is stopped
	cmd = exec.Command("sc", "query", serviceName)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Service stop verification failed: %v", err)
	}
	log.Println("Service stopped.")
}

// TestUninstallService verifies that the service can be uninstalled correctly.
func TestUninstallService(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "uninstall")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to uninstall service: %v", err)
	}

	// Verify uninstallation
	cmd = exec.Command("sc", "query", serviceName)
	err = cmd.Run()
	if err == nil {
		t.Fatalf("Service was not uninstalled successfully.")
	}
	log.Println("Service uninstalled.")
}

// TestServiceWithPrivileges tests the service behavior with different privilege levels.
func TestServiceWithPrivileges(t *testing.T) {
	privileges := []string{"medium", "high", "system"}

	for _, privilege := range privileges {
		t.Run(privilege, func(t *testing.T) {
			err := exec.Command("cmd.exe", "/C", "set PRIVILEGE_LEVEL="+privilege+" && go run main.go start").Run()
			if err != nil {
				t.Fatalf("Failed to start service with %s privilege: %v", privilege, err)
			}
			log.Printf("Service started with %s privilege.", privilege)

			// Wait for a bit to allow the service to start
			time.Sleep(10 * time.Second)

			err = exec.Command("cmd.exe", "/C", "set PRIVILEGE_LEVEL="+privilege+" && go run main.go stop").Run()
			if err != nil {
				t.Fatalf("Failed to stop service with %s privilege: %v", privilege, err)
			}
			log.Printf("Service stopped with %s privilege.", privilege)
		})
	}
}
