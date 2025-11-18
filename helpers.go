package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func ensureVaultDir() {
	vaultFilePath := getVaultFilePath()
	vaultDirPath := filepath.Dir(vaultFilePath)
	os.MkdirAll(vaultDirPath, 0700)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getVaultFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, vaultDir, vaultFile)
}

func copyToClipboard(text string) {
	var cmd *exec.Cmd
	switch {
	case runtime.GOOS == "darwin":
		cmd = exec.Command("pbcopy")
	case runtime.GOOS == "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case runtime.GOOS == "windows":
		cmd = exec.Command("clip")
	default:
		fmt.Println("Clipboard not supported on this OS.")
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("failed to access clipboard: ", err)
	}
	// Use Go to write to stdin
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start clipboard command:", err)
		stdin.Close()
		return
	}

	// Write the text
	if _, err := stdin.Write([]byte(text)); err != nil {
		fmt.Println("Failed to write to clipboard:", err)
	}
	stdin.Close()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Clipboard command failed:", err)
		return
	}

	fmt.Printf("Password copied to clipboard. Will clear in %.0f seconds.\n", clipboardDelay.Seconds())

	// Clear clipboard after delay
	time.AfterFunc(clipboardDelay, func() {
		clearCmd := exec.Command(cmd.Path, cmd.Args[1:]...) // reuse same command
		clearStdin, _ := clearCmd.StdinPipe()
		if err := clearCmd.Start(); err != nil {
			fmt.Println("Failed to start clipboard command:", err)
			return
		}
		spaces := strings.Repeat(" ", len(text))
		clearStdin.Write([]byte(spaces))
		clearStdin.Close()
		clearCmd.Wait()
	})
}
