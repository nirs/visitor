// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
)

func startService() {
	isService, err := svc.IsWindowsService()
	if err != nil || !isService {
		return // Not running as a Windows service.
	}

	// XXX Use %AppData%/visitor?
	execPath, err := os.Executable()
	if err == nil {
		_ = os.Chdir(filepath.Dir(execPath))
	}

	go func() {
		_ = svc.Run("", handler{})
		os.Exit(0)
	}()
}

type handler struct{}

func (handler) Execute(args []string, request <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	status <- svc.Status{
		State: svc.StartPending,
	}

	status <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	for {
		req := <-request
		switch req.Cmd {
		case svc.Interrogate:
			status <- req.CurrentStatus
		case svc.Stop, svc.Shutdown:
			status <- svc.Status{State: svc.StopPending}
			return false, 0
		}
	}
}
