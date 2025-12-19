package main

import (
	"bufio"
	"crypto/aes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	vaultDir           = ".vault"
	vaultFile          = "vault.json.aes"
	saltSize           = 32
	nonceSize          = 12
	argonTime          = 4
	argonMemory        = 64 * 1024
	argonThreads       = 2
	argonKeyLen        = 32
	scryptN            = 1 << 18
	scryptR            = 8
	scryptP            = 1
	MasterPasswordLen  = 12
	clipboardDelay     = 30 * time.Second
	recoveryQuestionsN = 3
	recoverySaltSize   = 32
)

func main() {
	ensureVaultDir()
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "init":
		initVault()
	case "add":
		addEntry()
	case "list":
		listEntries()
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Missing site name.")
			return
		}
		getEntry(os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: securepass delete <site>")
			return
		}
		deleteEntry(os.Args[2])
	case "change-master":
		changeMasterPassword()
	case "recover":
		recoverVault()
	default:
		printHelp()

	}

}

func initVault() {
	homeDir, _ := os.UserHomeDir()
	vaultFilePath := filepath.Join(homeDir, vaultDir, vaultFile)
	if fileExists(vaultFilePath) {
		fmt.Println("Vault already exists. Use 'change-master' to rotate key.")
		return
	}

	masterPass := readMasterPassword("Set master password: ")
	confirm := readMasterPassword("confirm master password: ")

	if masterPass != confirm {
		fmt.Println("Passwords do not match.")
		return
	}

	fmt.Println("Recovery Setup ... ")
	recoveryAnswers := askRecoveryQuestions()

	salt := generateSalt()
	nonce := generateNonce()
	key := deriveKeyArgon2(masterPass, salt)
	vault := Vault{Entries: []Entry{}}
	encrypted := encryptVault(vault, key, nonce)

	err := os.MkdirAll(filepath.Dir(vaultFilePath), 0700)
	if err != nil {
		panic(err)
	}
	data := append(salt, nonce...)
	data = append(data, encrypted...)
	err = os.WriteFile(vaultFilePath, data, 0600)
	if err != nil {
		panic(err)
	}

	err = appendRecoveryData(vaultFilePath, RecoveryData{Answers: recoveryAnswers})
	if err != nil {
		panic(err)
	}

	fmt.Println("Vault initialized securely at: ", vaultFile)
}

func recoverVault() {
	fmt.Println("Vault recovery initiated.")
	data, err := os.ReadFile(getVaultFilePath() + ".recovery")
	if err != nil {
		fmt.Println("Failed to read recovery data file:", err)
		os.Exit(1)
	}

	var recoveryData RecoveryData
	err = json.Unmarshal(data, &recoveryData)
	if err != nil {
		fmt.Println("Corrupted recovery data file:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	correctAnswers := 0

	for _, qa := range recoveryData.Answers {
		fmt.Println(qa.Question)
		answer := readLine(scanner, "Answer: ")
		normalizedAnswer := strings.ToLower(answer)
		hash := argon2.IDKey([]byte(normalizedAnswer), qa.Salt, argonTime, argonMemory, argonThreads, argonKeyLen)

		if compareHashes(hash, qa.Hash) {
			correctAnswers++
		} else {
			fmt.Println("Incorrect answer.")
		}
	}

	if correctAnswers < recoveryQuestionsN {
		fmt.Println("Recovery failed. Not enough correct answers.")
		os.Exit(1)
	}

	newMasterPass := readMasterPassword("Set new master password: ")
	confirm := readMasterPassword("Confirm new master password: ")

	if newMasterPass != confirm {
		fmt.Println("Passwords do not match.")
		return
	}

	vault, _, _ := loadVault()
	salt := generateSalt()
	nonce := generateNonce()
	key := deriveKeyArgon2(newMasterPass, salt)
	encrypted := encryptVault(vault, key, nonce)

	fullData := append(salt, nonce...)
	fullData = append(fullData, encrypted...)
	err = os.WriteFile(getVaultFilePath(), fullData, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Vault recovered and master password reset.")
}

func loadVault() (Vault, []byte, []byte) {
	data, err := os.ReadFile(getVaultFilePath())
	if err != nil {
		fmt.Println("No vault found. Run 'securepass init' first.")
		os.Exit(1)
	}

	if len(data) < saltSize+nonceSize+aes.BlockSize {
		fmt.Println("Corrupted vault.")
		os.Exit(1)
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	encryptedData := data[saltSize+nonceSize:]

	masterPass := readMasterPassword("Master password: ")
	key := deriveKeyArgon2(masterPass, salt)

	plaintext, err := decrypt(encryptedData, key, nonce)
	if err != nil {
		// Fallback to scrypt (for legacy)
		key = deriveKeyScrypt(masterPass, salt)
		plaintext, err = decrypt(encryptedData, key, nonce)
		if err != nil {
			fmt.Println("Incorrect password or corrupted vault.", err)
			os.Exit(1)
		}
	}

	var vault Vault
	json.Unmarshal(plaintext, &vault)
	return vault, key, nonce
}

func saveVault(vault Vault, key, nonce []byte) {
	encrypted := encryptVault(vault, key, nonce)
	salt := make([]byte, saltSize)
	// We need to read the salt from the existing vault file to re-encrypt
	data, err := os.ReadFile(getVaultFilePath())
	if err == nil {
		copy(salt, data[:saltSize])
	}
	fullData := append(salt, nonce...)
	fullData = append(fullData, encrypted...)
	os.WriteFile(getVaultFilePath(), fullData, 0600)
}

func addEntry() {
	vault, key, nonce := loadVault()

	var site, username, password string
	fmt.Print("Site: ")
	fmt.Scanln(&site)
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	fmt.Print("Password: ")
	fmt.Scanln(&password)
	vault.Entries = append(vault.Entries, Entry{site, username, password, ""})
	saveVault(vault, key, nonce)
	fmt.Println("Entry added.")
}

func listEntries() {
	vault, _, _ := loadVault()
	if len(vault.Entries) == 0 {
		fmt.Println("No entries.")
		return
	}
	fmt.Printf("%-25s %-20s %s\n", "Site", "Username", "Note")
	fmt.Println(strings.Repeat("-", 70))
	for _, e := range vault.Entries {
		note := e.Note
		if len(note) > 30 {
			note = note[:27] + "..."
		}
		fmt.Printf("%-25s %-20s %s\n", e.Site, e.Username, note)
	}
}

func getEntry(site string) {
	vault, _, _ := loadVault()
	for _, e := range vault.Entries {
		if strings.EqualFold(e.Site, site) {
			fmt.Printf("Site: %s\nUsername: %s\nPassword: %s\n", e.Site, e.Username, e.Password)
			if e.Note != "" {
				fmt.Printf("Note: %s\n", e.Note)
			}
			copyToClipboard(e.Password)
			return
		}
	}
	fmt.Println("Site not found.")
}

func deleteEntry(site string) {
	vault, key, nonce := loadVault()
	for i, e := range vault.Entries {
		if strings.EqualFold(e.Site, site) {
			vault.Entries = append(vault.Entries[:i], vault.Entries[i+1:]...)
			saveVault(vault, key, nonce)
			fmt.Println("Entry deleted.")
			return
		}
	}
	fmt.Println("Site not found.")
}

func changeMasterPassword() {
	if !fileExists(vaultFile) {
		fmt.Println("No vault to re-encrypt.")
		return
	}

	data, _ := os.ReadFile(vaultFile)
	salt := data[:saltSize]
	ciphertext := data[saltSize:]
	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	oldPass := readMasterPassword("Current master password: ")
	key := deriveKeyArgon2(oldPass, salt)
	plaintext, err := decrypt(encryptedData, key, nonce)
	if err != nil {
		key = deriveKeyScrypt(oldPass, salt)
		plaintext, err = decrypt(encryptedData, key, nonce)
		if err != nil {
			fmt.Println("Wrong password.")
			os.Exit(1)
		}
	}

	var vault Vault
	json.Unmarshal(plaintext, &vault)

	newPass := readMasterPassword("New master password: ")
	confirm := readMasterPassword("Confirm: ")
	if newPass != confirm {
		fmt.Println("Mismatch.")
		os.Exit(1)
	}

	newSalt := generateSalt()
	newKey := deriveKeyArgon2(newPass, newSalt)
	newNonce := generateNonce()
	newEncrypted := encryptVault(vault, newKey, newNonce)

	fullData := append(newSalt, append(newNonce, newEncrypted...)...)
	os.WriteFile(vaultFile, fullData, 0600)
	fmt.Println("Master password changed. Vault re-encrypted.")
}
