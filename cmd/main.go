package main

import (
	"fmt"
	"log"
	"net"
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
		return
	}

	targetIP := os.Getenv("TARGET_IP")
	if targetIP == "" {
		log.Fatal("TARGET_IP not set")
	}
	duration := os.Getenv("DURATION")
	if duration == "" {
		duration = "10"
	}

	parallel := os.Getenv("PARALLEL")
	if parallel == "" {
		parallel = "5"
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to list interfaces: %v", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Skipping interface %s: cannot get addresses: %v", iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil || ip.To4() == nil {
				continue // skip non-IPv4 or malformed
			}

			fmt.Printf("Testing bandwidth from %s (%s) to %s for %s seconds\n", iface.Name, ip.String(), targetIP, duration)
			cmd := exec.Command("iperf3", "-B", ip.String(), "-c", targetIP, "-t", duration, "-P", parallel)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Printf("iperf3 client on %s failed: %v", iface.Name, err)
			}
		}
	}
}
