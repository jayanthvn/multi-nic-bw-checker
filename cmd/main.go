package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	mode := os.Getenv("IPERF_MODE")
	if mode == "server" {
		fmt.Println("Starting iperf3 server with verbose logging...")
		cmd := exec.Command("iperf3", "-s", "--verbose")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("iperf3 server failed: %v", err)
		}
		// iperf3 should block and keep the pod running.
		return
	}

	ip := os.Getenv("TARGET_IP")
	if ip == "" {
		log.Fatal("TARGET_IP not set")
	}
	duration := os.Getenv("DURATION")
	if duration == "" {
		duration = "10"
	}

	fmt.Printf("Testing bandwidth to %s for %s seconds\n", ip, duration)
	cmd := exec.Command("iperf3", "-c", ip, "-t", duration)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("iperf3 client failed: %v", err)
	}
}
