package main

import (
	"log"
	"os/exec"
	"time"
)

func main() {
	// Run Ansible task and measure execution time
	executionTime := measureExecutionTime(runAnsibleTask)

	log.Printf("Ansible task executed in %s\n", executionTime)
}

func runAnsibleTask() {
	// Execute your Ansible task command
	cmd := exec.Command("ansible-playbook", "your-playbook.yaml")

	// Start the timer
	startTime := time.Now()

	// Run the command
	err := cmd.Run()

	// Calculate the execution time
	executionTime := time.Since(startTime)

	// Handle any error that occurred during command execution
	if err != nil {
		log.Printf("Error executing Ansible task: %v\n", err)
		return
	}

	log.Println("Ansible task executed successfully.")
}

func measureExecutionTime(taskFunc func()) time.Duration {
	// Start the timer
	startTime := time.Now()

	// Execute the task
	taskFunc()

	// Calculate the execution time
	executionTime := time.Since(startTime)

	return executionTime
}
