package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func memoryStress() {
	fileSizeMB := os.Getenv("DD_MB_SIZE")
	if fileSizeMB == "" {
		fileSizeMB = "60000"
	}
	fmt.Printf("File size set to %s MB.\n", fileSizeMB)

	sleepTime := os.Getenv("SLEEP_SEC")
	if sleepTime == "" {
		sleepTime = "60"
	}
	fmt.Printf("Sleep set to %s MB.\n", sleepTime)

	// 1. Directory creation with check if it already exists
	dir := "/tmp/pezhang"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, so create it
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}

	// 2. Mount tmpfs to /tmp/pezhang
	err := syscall.Mount("tmpfs", dir, "tmpfs", 0, "")
	if err != nil {
		// Check if already mounted (mount syscall can fail if it's already mounted)
		fmt.Printf("Error mounting tmpfs: %v\n", err)
		return
	}

	// 3. Generate a random file name
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFileName := fmt.Sprintf("%s/%d", dir, rng.Intn(1000000))

	// Convert fileSizeMB to an integer to validate input
	size, err := strconv.Atoi(fileSizeMB)
	if err != nil {
		fmt.Printf("Invalid DD_MB_SIZE value: %v using default 60000\n", err)
		size = 60000
	}

	// Create file with dd command using random file name
	cmd := exec.Command("dd", "if=/dev/zero", fmt.Sprintf("of=%s", randomFileName), "bs=1M", fmt.Sprintf("count=%d", size))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running dd command: %v\n", err)
		return
	}
	fmt.Printf("Created %sMB file at %s.\n", fileSizeMB, randomFileName)

	// Convert fileSizeMB to an integer to validate input
	seconds, err := strconv.Atoi(sleepTime)
	if err != nil {
		fmt.Printf("Invalid SLEEP_SEC value: %v using default 60\n", err)
		seconds = 60
	}

	if sleepTime == "0" {
		return
	}

	fmt.Println("Sleeping for %s sec...", sleepTime)
	time.Sleep(time.Duration(seconds) * time.Second)

	// 6. Delete the file after sleep
	err = os.Remove(randomFileName)
	if err != nil {
		fmt.Printf("Error deleting file %s: %v\n", randomFileName, err)
		return
	}
	fmt.Printf("File %s deleted, memory restored.\n", randomFileName)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := os.Getenv("HOSTNAME")
	if len(response) == 0 {
		response = "Hello OpenShift!"
	} else {
		response = "Service handled by pod " + response
	}

	// Echo back the port the request was received on
	// via a "request-port" header.
	addr := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		w.Header().Set("x-request-port", strconv.Itoa(tcpAddr.Port))
	}

	fmt.Fprintln(w, response)
	fmt.Println("Servicing request.")

	memoryStress()
}

func listenAndServe(port string) {
	fmt.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	// memoryStress()
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	port = os.Getenv("SECOND_PORT")
	if len(port) == 0 {
		port = "8888"
	}
	go listenAndServe(port)

	select {}
}
