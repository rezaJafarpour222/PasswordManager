package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
func SaveMasterKey(filePath string) error {
	password, err := encryption.GenerateRandomPassword(32)
	if err != nil {
		return err
	}
	key, err := base64.RawURLEncoding.DecodeString(password)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, key, 0600)
}
func LoadMasterKey(filePath string) (string, error) {
	key, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("master key not found")
	}

	if len(key) != 32 {
		return "", fmt.Errorf("invalid master key length: got %d bytes, want 32", len(key))
	}

	return base64.RawURLEncoding.EncodeToString(key), nil
}

func ExportVault(exportPath, vaultPath string) error {

	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return err
	}
	os.WriteFile(exportPath+"/Vault.vault", data, 0600)
	fmt.Println("Vault added to: ", exportPath)
	return nil
}
func ExportMasterKey(exportPath, masterKeyPath string) error {

	data, err := os.ReadFile(masterKeyPath)
	if err != nil {
		return err
	}
	err = os.WriteFile(exportPath+"/Master.key", data, 0600)
	if err != nil {
		return err
	}
	fmt.Println("Master added to: ", exportPath)
	return nil
}
