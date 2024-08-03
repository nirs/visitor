// Simplified version of https://github.com/caddyserver/caddy/blob/master/service_windows.go

package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
)

func init() {
	isService, err := svc.IsWindowsService()
	if err != nil || !isService {
		return
	}

	// Windows services always start in the system32 directory, try to
	// switch into the directory where the executable is.
	execPath, err := os.Executable()
	if err == nil {
		_ = os.Chdir(filepath.Dir(execPath))
	}

	go func() {
		_ = svc.Run("", runner{})

		// XXX Do graceful shutdown.
		os.Exit(0)
	}()
}

type runner struct{}

func (runner) Execute(args []string, request <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	// XXX Report svc.StartPending and report svc.Running when the web server
	// is listening.
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
