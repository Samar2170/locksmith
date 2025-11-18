package main

import "fmt"

func printHelp() {
	fmt.Println(`
SecurePass - Local Password Manager (Go)

Commands:
  init             Create new vault
  add              Add new password
  list             List all sites
  get <site>       Show & copy password
  delete <site>    Remove entry
  change-master    Rotate master password

All data encrypted locally with Argon2id + AES-256-GCM.
No cloud. No logs. No telemetry.
    `)
}
