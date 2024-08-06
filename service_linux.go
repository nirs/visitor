// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"net"
	"os"
)

const sdNotifyReady = "READY=1"

func startService() {
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket == "" {
		return // Not running as systemd service.
	}

	addr := &net.UnixAddr{Name: notifySocket, Net: "unixgram"}
	conn, err := net.DialUnix(addr.Net, nil, addr)
	if err != nil {
		log.Print("Cannot connecting to systemd socket %s: %s", notifySocket, err)
		return
	}

	defer conn.Close()
	if _, err := conn.Write([]byte(sdNotifyReady)); err != nil {
		log.Print("Cannot write to systemd socket %s: %s", notifySocket, err)
	}
}
