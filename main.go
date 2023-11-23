package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type workspace struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

var sessionName string

func init() {
	flag.StringVar(&sessionName, "name", "default", "session name in .sessionrc")
	_ = flag.Parsed()
}

func main() {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user:", err)
	}

	// Construct the full path to the ~/.sessionrc file
	sessionrcPath := filepath.Join(usr.HomeDir, ".sessionrc")

	f, err := os.Open(sessionrcPath)
	if err != nil {
		log.Fatalf("error occur when open file: %v", err)
	}
	defer f.Close()

	var workspaces map[string][]workspace
	data, _ := io.ReadAll(f)

	err = json.Unmarshal(data, &workspaces)
	if err != nil {
		log.Fatalf("error occur when open file: %v", err)
	}

	// create tmux workspace
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error occur when create tmux workspaces: %v", err)
	}

	for _, cfg := range workspaces[sessionName] {
		cmd := exec.Command(
			"tmux", "send-keys",
			"-t", sessionName,
			fmt.Sprintf("tmux new-window -n %s", cfg.Name),
			"C-m",
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Fatalf("error occur when create tmux workspaces: %v", err)
		}
	}

	for i, cfg := range workspaces[sessionName] {
		time.Sleep(time.Second)

		cmd := exec.Command(
			"tmux", "send-keys",
			"-t", fmt.Sprintf("%s:%d", sessionName, i+2),
			fmt.Sprintf("cd %s", cfg.Path),
			"C-m",
		)

		if err := cmd.Run(); err != nil {
			log.Fatalf("error occur when create tmux workspaces: %v", err)
		}
	}

	// if baseInd, err := getBaseWindowIndex(sessionName); err == nil {
	// 	time.Sleep(time.Second * 5)
	// 	// remove first tmux
	// 	cmd = exec.Command(
	// 		"tmux", "send-keys",
	// 		"-t", sessionName,
	// 		strings.Join([]string{"tmux kill-window -t", baseInd}, " "),
	// 		"C-m",
	// 	)
	//
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stdout
	//
	// 	if err := cmd.Run(); err != nil {
	// 		log.Fatalf("error occur when create tmux workspaces: %v", err)
	// 	}
	// }
}

func getBaseWindowIndex(sessionName string) (string, error) {
	// Run the tmux list-windows command
	cmd := exec.Command("tmux", "list-windows", "-t", sessionName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running tmux list-windows: %w", err)
	}

	// Parse the output to extract the index of the base window
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "(active)") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				return strings.Split(fields[0], "")[0], nil
			}
		}
	}

	return "", fmt.Errorf("base window not found in session '%s'", sessionName)
}
