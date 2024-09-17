package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/svc"
)

const (
	serviceName = "MyWindowsService"
)

type myService struct{}

func (m *myService) Execute(args []string, req <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	changes <- svc.Status{State: svc.StartPending}
	log.Println("Service started.")

	// Run command based on the privilege level
	privilege := os.Getenv("PRIVILEGE_LEVEL")
	log.Printf("Privilege Level: %s", privilege)

	err := runCommandWithPrivileges("cmd.exe", "/c", "echo Hello World")
	if err != nil {
		log.Printf("Error executing command: %v", err)
		changes <- svc.Status{State: svc.StopPending}
		return false, 1
	}

	changes <- svc.Status{State: svc.Running}

	for req := range req {
		switch req.Cmd {
		case svc.Stop:
			changes <- svc.Status{State: svc.StopPending}
			log.Println("Service stopping.")
			return false, 0
		default:
			log.Printf("Unexpected control request: %v", req)
		}
	}

	return true, 0
}

func runCommandWithPrivileges(command string, args ...string) error {
	privilege := os.Getenv("PRIVILEGE_LEVEL")
	var cmd *exec.Cmd

	switch privilege {
	case "medium":
		cmd = exec.Command(command, args...)
	case "high":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Start-Process %s -ArgumentList '%s' -Verb RunAs", command, argsToString(args)))
	case "system":
		return runAsSystem(command, args)
	default:
		return fmt.Errorf("unknown privilege level: %s", privilege)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	log.Printf("Command output:\n%s", output)
	return nil
}

func runAsSystem(command string, args []string) error {
	taskName := "TempTask"
	taskPath := `C:\Windows\System32\schtasks.exe`

	commandWithArgs := command + " " + argsToString(args)

	createTaskArgs := fmt.Sprintf("/Create /TN %s /TR \"%s\" /SC ONCE /ST 00:00 /RL HIGHEST /F", taskName, commandWithArgs)
	createCmd := exec.Command(taskPath, strings.Split(createTaskArgs, " ")...)
	if err := createCmd.Run(); err != nil {
		return fmt.Errorf("error creating scheduled task: %w", err)
	}

	runCmd := exec.Command(taskPath, "/Run", "/TN", taskName)
	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("error running scheduled task: %w", err)
	}

	deleteCmd := exec.Command(taskPath, "/Delete", "/TN", taskName, "/F")
	if err := deleteCmd.Run(); err != nil {
		return fmt.Errorf("error deleting scheduled task: %w", err)
	}

	return nil
}

func argsToString(args []string) string {
	return strings.Join(args, " ")
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			installService()
			return
		case "uninstall":
			uninstallService()
			return
		case "start":
			startService()
			return
		case "stop":
			stopService()
			return
		}
	}

	err := runService()
	if err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
}

func runService() error {
	err := svc.Run(serviceName, &myService{})
	if err != nil {
		return err
	}
	return nil
}

func installService() {
	binPath := `D:\Golang_Project\fourCore-project\windows-service\mywindowservice.exe`
	cmd := exec.Command("sc", "create", serviceName, "binPath=", binPath, "start=", "auto")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to install service: %v", err)
	}
	log.Println("Service installed.")
}

func startService() {
	cmd := exec.Command("sc", "start", serviceName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
	log.Println("Service started.")
}

func stopService() {
	cmd := exec.Command("sc", "stop", serviceName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to stop service: %v", err)
	}
	log.Println("Service stopped.")
}

func uninstallService() {
	cmd := exec.Command("sc", "delete", serviceName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to uninstall service: %v", err)
	}
	log.Println("Service uninstalled.")
}
