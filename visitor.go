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
	Time    time.Time `json:time`
	Count   uint64    `json:count`
	Address string    `json:address`
}

type State struct {
	Current *Visit
	Last    *Visit
}

func main() {
	http.HandleFunc("/", visit)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func visit(w http.ResponseWriter, r *http.Request) {
	state, err := update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h1>Welcome to %v</h1>", state.Current.Address)
	fmt.Fprintf(w, "<h2>You are visitor #%v</h2>", state.Current.Count)

	if state.Last.Address != "" {
		fmt.Fprintf(w, "<h2>Last visit on %v at %v</h2>", state.Last.Address, state.Last.Time)
	}
}

func update() (*State, error) {
	mutex.Lock()
	defer mutex.Unlock()

	last, err := readVisit()
	if err != nil {
		return nil, err
	}

	current := newVisit(last)
	if err := writeVisit(current); err != nil {
		return nil, err
	}

	return &State{Current: current, Last: last}, nil
}

func ipAddress() string {
	con, err := net.Dial("udp", "1.2.3.4:80")
	if err != nil {
		return "unknown"
	}
	return con.LocalAddr().(*net.UDPAddr).IP.String()
}

func newVisit(last *Visit) *Visit {
	var count uint64
	if last != nil {
		count = last.Count + 1
	}
	return &Visit{Time: time.Now(), Count: count, Address: ipAddress()}
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
