package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

type workspace struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func main() {
	f, err := os.Open("./workspace.json")
	if err != nil {
		log.Fatalf("error occur when open file: %v", err)
	}
	defer f.Close()

	var workspaces []workspace
	data, _ := io.ReadAll(f)

	err = json.Unmarshal(data, &workspaces)
	if err != nil {
		log.Fatalf("error occur when open file: %v", err)
	}

	for _, cfg := range workspaces {
		if err := os.Chdir(cfg.Path); err != nil {
			log.Fatalf("error occur when create tmux workspaces: %v", err)
		}

		cmd := exec.Command("tmux", "new-window", "-n", cfg.Name)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Fatalf("error occur when create tmux workspaces: %v", err)
		}
	}

}
