// Yanked this from absdev.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// Prints out text live while also capturing a copy
func captureAndPrintOutput(r io.ReadCloser, wg *sync.WaitGroup) string {
	defer wg.Done()

	output := ""
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		output += line + "\n"
		fmt.Println(line)
	}
	return output
}

func RunInteractiveCommand(name string, arg ...string) (stdout string, stderr string, err error) {
	cmd := NewCommandBuilder(name, arg...)
	return cmd.RunInteractive()
}

func RunInteractiveCommandNoPipes(name string, arg ...string) error {
	cmd := NewCommandBuilder(name, arg...)
	return cmd.RunInteractiveNoPipes()
}

// Do not automatically print output, but return it.
func SilentCommand(name string, arg ...string) (output string, err error) {
	cmd := NewCommandBuilder(name, arg...)
	return cmd.RunSilent()
}

// Get the returncode from an error (or lack thereof)
func GetExitStatus(err error) (exitStatus int) {
	if err == nil {
		return
	}

	exitErr, ok := err.(*exec.ExitError)
	if ok {
		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			exitStatus = status.ExitStatus()
		}
	}
	return
}

// Allows someone to more easily setup and execute a single command
type CommandBuilder struct {
	Cmd *exec.Cmd
}

func NewCommandBuilder(name string, arg ...string) *CommandBuilder {
	return &CommandBuilder{exec.Command(name, arg...)}
}

// Add a key to the environment in additional to the current environment
func (c *CommandBuilder) AddEnv(key, val string) {
	if c.Cmd.Env == nil {
		c.Cmd.Env = os.Environ()
	}
	c.Cmd.Env = append(c.Cmd.Env, key+"="+val)
}

func (c *CommandBuilder) RunInteractive() (stdout string, stderr string, err error) {
	c.Cmd.Stdin = os.Stdin
	stdoutPipe, err := c.Cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderrPipe, err := c.Cmd.StderrPipe()
	if err != nil {
		stdoutPipe.Close()
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		stdout = captureAndPrintOutput(stdoutPipe, wg)
	}()
	go func() {
		stderr = captureAndPrintOutput(stderrPipe, wg)
	}()

	err = c.Cmd.Start()
	if err != nil {
		return
	}
	wg.Wait()
	return stdout, stderr, c.Cmd.Wait()
}

// Some commands have different output that is not friendly with pipes. This avoids capturing the output
// and instead links the output directly to stdout/stderr.
func (c *CommandBuilder) RunInteractiveNoPipes() error {
	c.Cmd.Stdin = os.Stdin
	c.Cmd.Stdout = os.Stdout
	c.Cmd.Stderr = os.Stderr
	return c.Cmd.Run()
}

func (c *CommandBuilder) RunSilent() (output string, err error) {
	var outputBytes []byte
	outputBytes, err = c.Cmd.CombinedOutput()
	if len(outputBytes) > 0 {
		output = string(outputBytes)
	}
	return
}
