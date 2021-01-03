package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func ExecuteCommand(command string, timeout int, cwd string, env []string) (string, string, int) {
	return "", "not support running cmd in windows", -1
}

func ExecuteCmd(command string, timeout int, cwd string, env []string) (string, string, int) {
	rc := -1
	if f, err := os.Stat(cwd); err != nil || !f.IsDir() {
		msg := fmt.Sprintf("chdir failed, dir %s not exists!", cwd)
		return "", msg, -500
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Dir = cwd

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return string(buf.Bytes()), err.Error(), -500
	}

	err := cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				rc = status.ExitStatus()
			}
		}
		return string(buf.Bytes()), err.Error(), rc
	}

	return string(buf.Bytes()), "", 0
}
