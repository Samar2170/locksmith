package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	recoveryQuestionsN = 2
	recoverySaltSize   = 32
)

type RecoveryAnswer struct {
	Question string `json:"q"`
	Salt     []byte `json:"s"` // per-question salt
	Hash     []byte `json:"h"` // Argon2id(answer + salt)
}

type RecoveryData struct {
	Answers []RecoveryAnswer `json:"answers"`
}

func askRecoveryQuestions() []RecoveryAnswer {
	fmt.Printf("Set up %d recovery questions (case-insensitive, trimmed):\n", recoveryQuestionsN)
	answers := make([]RecoveryAnswer, recoveryQuestionsN)

	for i := range recoveryQuestionsN {
		var q, a string
		fmt.Printf("Question %d: ", i+1)
		fmt.Scanln(&q)

		var dummy string
		fmt.Scanln(&dummy)

		fmt.Print("Answer: ")
		fmt.Scanln(&a)
		a = strings.TrimSpace(strings.ToLower(a))
		salt := generateSalt()
		hash := argon2.IDKey([]byte(a), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
		answers[i] = RecoveryAnswer{
			Question: q,
			Salt:     salt,
			Hash:     hash,
		}
	}
	return answers
}

func appendRecoveryData(path string, data RecoveryData) error {
	b, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(path+".recovery", b, 0600)
}
