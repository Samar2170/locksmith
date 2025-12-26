package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

type RecoveryAnswer struct {
	Question string `json:"q"`
	Salt     []byte `json:"s"` // per-question salt
	Hash     []byte `json:"h"` // Argon2id(answer + salt)
}

type RecoveryData struct {
	Answers           []RecoveryAnswer `json:"answers"`
	EncryptedVaultKey []byte           `json:"encrypted_vault_key"`
	KeyNonce          []byte           `json:"key_nonce"`
	RecoveryKeySalt   []byte           `json:"recovery_key_salt"`
}

func askRecoveryQuestions() ([]RecoveryAnswer, []string) {

	fmt.Printf("Set up %d recovery questions (answers are case-insensitive and trimmed):\n\n", recoveryQuestionsN)

	answers := make([]RecoveryAnswer, recoveryQuestionsN)
	rawAnswers := make([]string, recoveryQuestionsN)

	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i < recoveryQuestionsN; i++ {
		question := readLine(scanner, fmt.Sprintf("Question %d: ", i+1))
		for question == "" {
			fmt.Println("Question cannot be empty.")
			question = readLine(scanner, fmt.Sprintf("Question %d: ", i+1))
		}

		answer := readLine(scanner, "Answer: ")
		for answer == "" {
			fmt.Println("Answer cannot be empty.")
			answer = readLine(scanner, "Answer: ")
		}

		// Normalize and hash the answer
		normalizedAnswer := strings.ToLower(answer)
		salt := generateSalt()
		hash := argon2.IDKey([]byte(normalizedAnswer), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

		answers[i] = RecoveryAnswer{
			Question: question,
			Salt:     salt,
			Hash:     hash,
		}
		rawAnswers[i] = answer

		fmt.Println("Saved.\n")
	}

	return answers, rawAnswers

}

func appendRecoveryData(path string, data RecoveryData) error {
	b, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(path+".recovery", b, 0600)
}

func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
