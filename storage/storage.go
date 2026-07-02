package storage

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"pass/encryption"
)

func SaveVault(v encryption.Vault, password, path string) error {
	salt, err := encryption.GenerateSalt()
	if err != nil {
		return err
	}

	key := encryption.DeriveKey(password, salt)

	plain, err := json.Marshal(v)
	if err != nil {
		return err
	}

	nonce, ciphertext, err := encryption.Encrypt(key, plain)
	if err != nil {
		return err
	}

	file := encryption.VaultFile{
		Salt:       base64.StdEncoding.EncodeToString(salt),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		CipherText: base64.StdEncoding.EncodeToString(ciphertext),
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadVault(password, path string) (encryption.Vault, error) {
	var v encryption.Vault
	var file encryption.VaultFile

	data, err := os.ReadFile(path)
	if err != nil {
		return v, err
	}

	json.Unmarshal(data, &file)

	salt, _ := base64.StdEncoding.DecodeString(file.Salt)
	nonce, _ := base64.StdEncoding.DecodeString(file.Nonce)
	ciphertext, _ := base64.StdEncoding.DecodeString(file.CipherText)

	key := encryption.DeriveKey(password, salt)

	plain, err := encryption.Decrypt(key, nonce, ciphertext)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(plain, &v)
	return v, err
}
