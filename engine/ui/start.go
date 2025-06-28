package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)


func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // "linux", "freebsd", etc.
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Run()
}

func StartReactApp(projectRoot string) error {
	// Change to the ui directory
	uiDir := filepath.Join(projectRoot, "ui")
	
	// Install dependencies if node_modules doesn't exist
	if _, err := os.Stat(filepath.Join(uiDir, "node_modules")); os.IsNotExist(err) {
		installCmd := exec.Command("npm", "install")
		installCmd.Dir = uiDir
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install dependencies: %v", err)
		}
	}

	// Start the React development server
	npmCmd := exec.Command("npm", "start")
	npmCmd.Dir = uiDir
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr
	// Set PORT=3000 environment variable
	npmCmd.Env = append(os.Environ(), "PORT=3000", "BROWSER=none") // BROWSER=none prevents npm from opening the browser
	
	if err := npmCmd.Start(); err != nil {
		return fmt.Errorf("failed to start React app: %v", err)
	}

	// Wait a bit for the server to start
	time.Sleep(3 * time.Second)
	
	// Open the browser
	if err := openBrowser("http://localhost:3000"); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}

	// Don't wait for the process to finish - let it run in background
	return nil
}
