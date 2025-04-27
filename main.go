package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	// Find all service directories (containing main.go files) under cmd/
	cmdDir := filepath.Join(rootDir, "cmd")
	services, err := findServices(cmdDir)
	if err != nil {
		fmt.Printf("Error finding services: %v\n", err)
		os.Exit(1)
	}
	if len(services) == 0 {
		fmt.Println("No services found in cmd/ directory")
		os.Exit(1)
	}
	fmt.Printf("Found %d services to start: %v\n", len(services), services)
	// Create a wait group to manage goroutines
	var wg sync.WaitGroup
	// Start each service in a separate goroutine
	for _, service := range services {
		wg.Add(1)
		go startService(service, &wg)
	}
	time.Sleep(20 * time.Second)
	fmt.Println("All services started. Press Ctrl+C to stop all services.")
	wg.Wait()
}

func findServices(cmdDir string) ([]string, error) {
	var services []string
	// Walk through the cmd directory
	err := filepath.Walk(cmdDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If this is a main.go file, add its directory to services
		if !info.IsDir() && info.Name() == "main.go" {
			serviceName := filepath.Base(filepath.Dir(path))
			services = append(services, serviceName)
		}

		return nil
	})
	return services, err
}

func startService(serviceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting service: %s\n", serviceName)
	// Get the path to the service's main.go file
	rootDir, _ := os.Getwd()
	servicePath := filepath.Join(rootDir, "cmd", serviceName)
	// Determine command to run based on OS
	var cmd *exec.Cmd
	cmd = exec.Command("go", "run", filepath.Join(servicePath, "main.go"))
	// Set the working directory for the command
	cmd.Dir = rootDir
	// Connect standard outputs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run the command
	err := cmd.Start()
	fmt.Printf("Complete the startup service: %s\n", serviceName)
	if err != nil {
		fmt.Printf("Error starting service %s: %v\n", serviceName, err)
		return
	}
	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Service %s exited with error: %v\n", serviceName, err)
	} else {
		fmt.Printf("Service %s completed successfully\n", serviceName)
	}
}
