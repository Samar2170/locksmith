package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

func generateSalt() []byte {
	salt := make([]byte, saltSize)
	rand.Read(salt)
	return salt
}

func generateNonce() []byte {
	nonce := make([]byte, nonceSize)
	rand.Read(nonce)
	return nonce
}

func deriveKeyArgon2(pass string, salt []byte) []byte {
	return argon2.IDKey([]byte(pass), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}

func encrypt(data, key, nonce []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return aesgcm.Seal(nil, nonce, data, nil)
}

func encryptVault(vault Vault, key, nonce []byte) []byte {
	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return aesgcm.Seal(nil, nonce, data, nil)
}
func decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func deriveKeyScrypt(pass string, salt []byte) []byte {
	key, _ := scrypt.Key([]byte(pass), salt, scryptN, scryptR, scryptP, argonKeyLen)
	return key
}
