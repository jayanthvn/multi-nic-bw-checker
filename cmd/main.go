package main

import (
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	serverMode := os.Getenv("IPERF_MODE") == "server"
	target := os.Getenv("TARGET_IP")
	duration := os.Getenv("DURATION")
	if duration == "" {
		duration = "10"
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Failed to get interfaces: %v", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			if ip == nil || ip.To4() == nil {
				continue
			}

			if serverMode {
				cmd := exec.Command("iperf3", "-s")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				log.Printf("Starting iperf3 server on interface: %s", iface.Name)
				_ = cmd.Run()
				return
			} else if target != "" {
				log.Printf("Running iperf3 test from %s to %s", ip.String(), target)
				cmd := exec.Command("iperf3", "-c", target, "-B", ip.String(), "-t", duration)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			}
		}
	}
}

