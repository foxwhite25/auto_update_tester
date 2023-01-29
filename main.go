package main

import (
	"os/exec"
	"strings"
)

func ExecuteWithOutput(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func RestartSelf() error {
	args := []string{"run", "auto_update_tester"}
	_, err := ExecuteWithOutput("go", args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	println("Updated Version 1")

	pullResult, err := ExecuteWithOutput("git", []string{"pull"})
	if err != nil {
		panic(err)
	}

	if strings.Contains(pullResult, "Already up to date.") {
		println("Already up to date, exiting...")
		return
	}
	if strings.Contains(pullResult, "CONFLICT") {
		println("There is a conflict, exiting...")
		return
	}

	println("Pull success, restart self.")
	err = RestartSelf()
	if err != nil {
		panic(err)
	}
}
