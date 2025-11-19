
generate_salt -> generate_key(masterPass,salt) -> aes.NewCipher(key) -> cipher.NewGCM(block) -> aesgcm.Seal