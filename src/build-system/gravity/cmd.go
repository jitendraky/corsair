package main

import (
  "os/exec"
  "syscall"
  "log"
)

func StartCommand(workingDirectory, shCommand string) *exec.Cmd {
  cmd := exec.Command("/bin/sh", "-c", shCommand)
	cmd.Dir = workingDirectory
  err := cmd.Start()
  if err != nil {
    log.Println("[Error] Gravity error occured during command: ", shCommand, ": ", err)
  }
	return cmd
}

func KillProcess(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	err := cmd.Wait()
	if err != nil {
		log.Println(err)
	}
	cmd = nil
}
