package lib

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
)

// Callback - shell output handler
type Callback func(line string)

// Shell - shell executor
func Shell(command string, callback Callback) {
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}

	scanner := bufio.NewScanner(out)

	if err = cmd.Start(); err != nil {
		log.Fatal("Unable to start command: ", err)
	}

	for scanner.Scan() {
		callback(scanner.Text())
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Unable to wait for process to finish: ", err)
	}
}

// Run - simple shell execution
func Run(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	result, err := cmd.Output()
	return string(result), err
}

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}
