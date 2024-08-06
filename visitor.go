// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const visitorFile = "visitor.json"

var mutex sync.Mutex

type Visit struct {
	Time  time.Time `json:"time"`
	Count uint64    `json:"count"`
	Host  string    `json:"host"`
}

type State struct {
	Current *Visit
	Last    *Visit
}

func main() {
	startService()
	http.HandleFunc("/", visit)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func visit(w http.ResponseWriter, r *http.Request) {
	state, err := update(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h1>Welcome to %v</h1>", state.Current.Host)
	fmt.Fprintf(w, "<h2>You are visitor #%v</h2>", state.Current.Count)

	if state.Last.Host != "" {
		fmt.Fprintf(w, "<h2>Last visit on %v at %v</h2>", state.Last.Host, state.Last.Time)
	}
}

func update(r *http.Request) (*State, error) {
	mutex.Lock()
	defer mutex.Unlock()

	last, err := readVisit()
	if err != nil {
		return nil, err
	}

	current := newVisit(last, r)
	if err := writeVisit(current); err != nil {
		return nil, err
	}

	return &State{Current: current, Last: last}, nil
}

func newVisit(last *Visit, r *http.Request) *Visit {
	var count uint64
	if last != nil {
		count = last.Count + 1
	}
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = ""
	}
	return &Visit{Time: time.Now(), Count: count, Host: host}
}

func readVisit() (*Visit, error) {
	visit := &Visit{}
	data, err := os.ReadFile(visitorFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return visit, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, visit); err != nil {
		return nil, err
	}
	return visit, nil
}

func writeVisit(visit *Visit) error {
	var err error
	data, err := json.Marshal(visit)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(visitorFile), filepath.Base(visitorFile)+".tmp")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
		}
	}()
	if _, err = tmp.Write(data); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil {
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmp.Name(), visitorFile); err != nil {
		return err
	}
	return nil
}
