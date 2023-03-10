package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func RunCommandWithArgs(command string, args []string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
	return syscall.Exec(command, args, os.Environ())
}

func ExecuteWithOutput(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func RestartSelf() error {
	args := []string{"run", "auto_update_tester"}
	err := RunCommandWithArgs("go", args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	println("Version 2")

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
