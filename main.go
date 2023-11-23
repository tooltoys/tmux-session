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
	flag.Parse()
}

func main() {
	// kill session
	exec.Command("tmux", "kill-session", "-t", sessionName).Run()

	session := sessions()

	for i, cfg := range session {
		if i == 0 {
			// create tmux session
			cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-n", cfg.Name)
			if err := cmd.Run(); err != nil {
				log.Fatalf("error occur when create tmux workspaces: %v", err)
			}
			continue
		}

		cmd := exec.Command(
			"tmux", "send-keys",
			"-t", sessionName,
			fmt.Sprintf("tmux new-window -n %s", cfg.Name),
			"C-m",
		)

		if err := cmd.Run(); err != nil {
			log.Fatalf("error occur when create tmux workspaces: %v", err)
		}
	}

	for i, cfg := range session {
		time.Sleep(time.Second)
		cmd := exec.Command(
			"tmux", "send-keys",
			"-t", fmt.Sprintf("%s:%d", sessionName, i+1),
			fmt.Sprintf("cd %s", cfg.Path),
			"C-m",
		)

		if err := cmd.Run(); err != nil {
			log.Fatalf("error occur when create tmux workspaces: %v", err)
		}
	}
}

func sessions() []workspace {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %v", err)
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

	return workspaces[sessionName]
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
