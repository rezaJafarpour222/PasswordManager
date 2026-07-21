package socket

import (
	"fmt"
	"os"
	"pass/encryption"
	"pass/storage"
	"strconv"
)

func Gen(service, username string, sizeStr string, v *encryption.Vault) (encryption.Entry, error) {

	size, err := strconv.Atoi(sizeStr)
	var entry encryption.Entry
	if err != nil {
		return entry, fmt.Errorf("Password length must be a number.")
	}
	if size > 33 {
		return entry, fmt.Errorf("Password length must be <=32.")
	}

	password, err := encryption.GenerateRandomPassword(size)
	if err != nil {
		return entry, err
	}
	entry = encryption.Entry{
		Service:  service,
		Username: username,
		Password: password,
	}

	v.AddEntry(entry)
	return entry, nil
}
func Add(service, username, password string, v *encryption.Vault) {

	entry := encryption.Entry{
		Service:  service,
		Username: username,
		Password: password,
	}
	v.AddEntry(entry)
}
func List(v *encryption.Vault) []encryption.Entry {
	return v.Entries
}

func DeleteByUserName(username string, v *encryption.Vault) error {
	err := v.DeleteEntryByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func DeleteByService(service string, v *encryption.Vault) error {
	err := v.DeleteEntryByService(service)
	if err != nil {
		return err
	}
	return nil
}

func Export(exportPath, currentPath string) error {

	err := storage.ExportVault(exportPath, currentPath)
	if err != nil {
		return err
	}
	err = storage.ExportMasterKey(exportPath, currentPath)
	if err != nil {
		return err
	}
	return nil
}

func Init(path string) error {
	vaultPath := path + "vault.Vault"
	masterKeyPath := path + "Master.key"
	_, err := os.Stat(vaultPath)
	if err == nil {
		return fmt.Errorf("Vault does exist")
	}
	_, err = os.Stat(masterKeyPath)
	if err == nil {
		return fmt.Errorf("Master key does exist")
	}
	err = storage.SaveMasterKey(masterKeyPath)
	if err != nil {
		return fmt.Errorf("Problem Loading Master key")
	}
	masterKey, err := storage.LoadMasterKey(vaultPath)
	if err != nil {
		return fmt.Errorf("Problem Loading Vault")
	}
	v := encryption.NewVault()
	err = storage.SaveVault(v, masterKey, vaultPath)
	if err != nil {
		return fmt.Errorf("Problem Saving Vault")
	}
	return nil
}
