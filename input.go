package main

import (
	"bufio"
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func readMasterPassword(prompt string) string {
	fmt.Print(prompt)
	bytePass, _ := term.ReadPassword(int(syscall.Stdin))
	if len(string(bytePass)) < MasterPasswordLen {
		fmt.Printf("Password must be at least %d characters long.", MasterPasswordLen)
		return readMasterPassword(prompt)
	}
	fmt.Println()
	return string(bytePass)
}

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
