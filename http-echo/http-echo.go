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

var MemoryTotalMb int64 = 0
var MemoryOccMb int64 = 0
var fileSizeMBInt int = 0
var memStressInt int = 0
var stop_allocating bool = false

func sleep(sleepTime string) int {
	if sleepTime == "0" {
		return 0
	}

	seconds := env_var_int(sleepTime, "sleep", 60)
	fmt.Printf("Sleeping for %d sec...\n", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)

	return seconds
}

func free_memory(randomFileName string) {
	err := os.Remove(randomFileName)
	if err != nil {
		fmt.Printf("Error deleting file %s: %v\n", randomFileName, err)
		return
	}
	fmt.Printf("File %s deleted, memory restored.\n", randomFileName)
}

func get_env_var(name string, empty string, def string) string {
	env := os.Getenv(name)
	if env == empty {
		env = def
	}
	return env
}

func env_var_int(name string, v string, def int) int {
	size, err := strconv.Atoi(name)
	if err != nil {
		fmt.Printf("Invalid %s value: %v using default %d\n", v, err, def)
		size = def
	}
	return size
}

func occupy_memory(randomFileName string) {
	if stop_allocating == false {
		MemoryOccMb += int64(fileSizeMBInt)
	}
	memory_occ := float64(MemoryOccMb * 100 / MemoryTotalMb)
	fmt.Printf("Memory occupied: %f%%\n", memory_occ)

	if stop_allocating {
		return
	}

	cmd := exec.Command("dd", "if=/dev/zero", fmt.Sprintf("of=%s", randomFileName), "bs=1M", fmt.Sprintf("count=%d", fileSizeMBInt))
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running dd command: %v\n", err)
		return
	}

	if memory_occ > float64(memStressInt) {
		if stop_allocating == false {
			/* Don't overload the server memory, don't increase memory anymore */
			fmt.Println("#### Stopping allocation to prevent server from going OOMKilled ####")
			stop_allocating = true
		}
		return
	}
	// fmt.Printf("Created %sMB file at %s.\n", fileSizeMBInt, randomFileName)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	dir := "/tmp/pezhang"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}

	err := syscall.Mount("tmpfs", dir, "tmpfs", 0, "")
	if err != nil {
		fmt.Printf("Error mounting tmpfs: %v\n", err)
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFileName := fmt.Sprintf("%s/%d", dir, rng.Intn(1000000))

	host := get_env_var("HOSTNAME", "", "???")
	response := "Service handled by pod " + host

	// Echo back the port the request was received on
	// via a "request-port" header.
	addr := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		w.Header().Set("x-request-port", strconv.Itoa(tcpAddr.Port))
	}

	fmt.Fprintln(w, response)
	fmt.Printf("---------------------------------\n")
	fmt.Printf("%s: processing request.\n", host)

	occupy_memory((randomFileName))
	// sleepTime := get_env_var("SLEEP_SEC", "", "60")
	// sl := sleep(sleepTime)
	// if sl > 0 {
	// 	free_memory(randomFileName)
	// }
}

func listenAndServe(port string) {
	fmt.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func parseMemoryString(memoryStr string) int64 {
	value, err := strconv.ParseInt(memoryStr, 10, 64)
	if err != nil {
		fmt.Printf("invalid memory format: %s\n", memoryStr)
		return 0
	}

	return value / (1024 * 1024)
}

func main() {
	// sleepTime := get_env_var("SLEEP_SEC", "", "60")
	// fmt.Printf("Sleep set to %s sec.\n", sleepTime)

	MemoryTotal := get_env_var("MEMORY_LIMITS", "", "???")
	MemoryTotalMb = parseMemoryString(MemoryTotal)
	fmt.Printf("Memory available in this pod is %d Mb.\n", MemoryTotalMb)

	fileSizeMB := get_env_var("DD_MB_SIZE", "", "60000")
	fileSizeMBInt = env_var_int(fileSizeMB, "DD_MB_SIZE", 60000)
	fmt.Printf("File size set to %s MB.\n", fileSizeMB)

	memStress := get_env_var("MAX_MEMORY_STRESS_PERC", "", "80")
	memStressInt = env_var_int(memStress, "MAX_MEMORY_STRESS_PERC", 80)
	fmt.Printf("Max memory stress set to %s%%.\n", memStress)

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
