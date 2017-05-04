package main

import (
	"bufio"
	"bytes"
	"os/exec"
)

func execCmd(command string, args ...string) *bytes.Buffer {
	ret := &bytes.Buffer{}
	retBuf := bufio.NewWriter(ret)
	cmd := exec.Command(command, args...)
	cmd.Stdout = retBuf
	cmd.Stderr = retBuf
	cmd.Run()
	return ret
}
