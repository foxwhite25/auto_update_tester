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
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	err = RunCommandWithArgs(self, args)
	if err == nil {
		os.Exit(0)
	}
	return err
}

func main() {
	println("Hello, world!")
	localHash, err := ExecuteWithOutput("git", []string{"rev-parse", "HEAD"})
	remoteHash, err := ExecuteWithOutput("git", []string{"ls-remote", "origin", "HEAD"})
	remoteHash = strings.Split(remoteHash, "\t")[0]
	if err != nil {
		panic(err)
	}
	if localHash != remoteHash {
		pullResult, err := ExecuteWithOutput("git", []string{"pull"})
		if err != nil {
			panic(err)
		}
		if strings.Contains(pullResult, "CONFLICT") {
			println("There is a conflict, ignore it.")
		} else {
			println("Pull success, restart self.")
			err = RestartSelf()
			if err != nil {
				panic(err)
			}
		}
	} else {
		println("No update detected, exiting...")
	}
}
