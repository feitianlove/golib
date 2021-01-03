package utils

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func ExecuteCommand(command string, timeout int, cwd string, env []string) (string, string, int) {
	rc := -1
	if f, err := os.Stat(cwd); err != nil || !f.IsDir() {
		msg := fmt.Sprintf("chdir failed, dir %s not exists!", cwd)
		log.Error(msg)
		return "", msg, -500
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Dir = cwd
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	errpipe, err := cmd.StderrPipe()
	if err != nil {
		return "", "cmd stderr pipe error", -500
	}
	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "cmd stdout pipe error", -500
	}
	err = cmd.Start()
	if err != nil {
		return "", "cmd start error", -500
	}

	copyOutDone := make(chan bool)
	copyErrDone := make(chan bool)
	var mout, merr bytes.Buffer
	go func() {
		_, err := io.Copy(&mout, outpipe)
		if err != nil {
			log.Error(err)
		}
		copyOutDone <- true
	}()

	go func() {
		_, err := io.Copy(&merr, errpipe)
		if err != nil {
			log.Error(err)
		}
		copyErrDone <- true
	}()

	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		rc = 500
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	case err := <-done:
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					rc = status.ExitStatus()
				}
			}
		} else {
			rc = 0
		}
		<-copyOutDone
		<-copyErrDone
	}

	return mout.String(), merr.String(), rc
}

func ExecuteCmd(command string, timeout int, cwd string, env []string) (string, string, int) {
	rc := -1
	if f, err := os.Stat(cwd); err != nil || !f.IsDir() {
		msg := fmt.Sprintf("chdir failed, dir %s not exists!", cwd)
		log.Error(msg)
		return "", msg, -500
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Dir = cwd
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
