package cli

import (
	"fmt"
	"os"
	"pass/encryption"
	"pass/storage"
	"strconv"
	"sync"
)

func (a *App) registerCommand(name, desc, example string) {
	a.Commands[name] = Command{
		Name:        name,
		Description: desc,
		Example:     example,
	}
}
func (a *App) Help() {

	fmt.Println("Usage:")
	fmt.Println("Command  Description")
	for _, cmd := range a.Commands {
		fmt.Println("#====================================================================#")
		fmt.Printf("%s  %s\n", cmd.Name, cmd.Description)
		fmt.Printf("Example -> %s\n", cmd.Example)
	}
}
func (a *App) Init() error {
	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)
	go Spinner(done, &wg, "Initializing the Vault and Master Key")
	_, err := os.Stat(a.VaultPath)
	if err == nil {
		return fmt.Errorf("vault does exist!")
	}
	err = storage.SaveMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}
	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}
	v := encryption.NewVault()
	err = storage.SaveVault(v, masterKey, a.VaultPath)
	if err != nil {
		return err
	}
	close(done)
	wg.Wait()
	fmt.Printf("%s and %s created.", a.VaultPath, a.MasterKeyPath)
	return nil
}

func (a *App) List() error {

	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)
	go Spinner(done, &wg, "Decrypting the Vault")
	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}
	v, err := storage.LoadVault(masterKey, a.VaultPath)
	if err != nil {
		return err
	}

	close(done)
	wg.Wait()
	if len(v.Entries) == 0 {
		fmt.Printf("Vault is empty\n")
		return nil
	}
	fmt.Println("#==================================#")
	for _, entry := range v.Entries {
		fmt.Println("Service: ", entry.Service)
		fmt.Println("Username: ", entry.Username)
		fmt.Println("Password: ", entry.Password)
		fmt.Println("#==================================#")
	}
	return nil

}

func (a *App) Add(service, username, password string) error {

	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)

	spinnerText := fmt.Sprintf("Adding '%s' to the vault", service)
	go Spinner(done, &wg, spinnerText)
	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}
	v, err := storage.LoadVault(masterKey, a.VaultPath)
	if err != nil {
		return err
	}

	v.AddEntry(encryption.Entry{
		Service:  service,
		Username: username,
		Password: password,
	})
	err = storage.SaveVault(v, masterKey, a.VaultPath)

	close(done)
	wg.Wait()

	if err != nil {
		return err
	}
	return nil
}

func (a *App) Gen(service, username string, sizeStr string) error {

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return fmt.Errorf("length must be a number")
	}

	randomPassword, err := encryption.GenerateRandomPassword(size)
	if err != nil {
		return err
	}
	m := make(map[string]string)
	m["Service"] = service
	m["Username"] = username
	m["Password"] = randomPassword
	if !Confirmation("gen", m) {
		return nil
	}

	err = a.Add(service, username, randomPassword)

	return nil
}

func (a *App) DeleteEntry(service string) error {

	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)
	spinnerText := fmt.Sprintf("Deleting '%s' from the vault", service)
	if !Confirmation("del", map[string]string{"Service": service}) {
		return nil
	}
	go Spinner(done, &wg, spinnerText)

	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}
	v, err := storage.LoadVault(masterKey, a.VaultPath)
	if err != nil {
		return err
	}
	err = v.DeleteEntry(service)
	if err != nil {
		return err
	}
	err = storage.SaveVault(v, masterKey, a.VaultPath)
	if err != nil {
		return err
	}

	close(done)
	wg.Wait()
	return nil
}

func (a *App) GetEntry(service string) error {

	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}
	v, err := storage.LoadVault(masterKey, a.VaultPath)
	if err != nil {
		return err
	}
	entries, err := v.GetEntry(service)
	if err != nil {
		return err
	}

	fmt.Println("#==================================#")
	for _, entry := range entries {
		fmt.Println("Service: ", entry.Service)
		fmt.Println("Username: ", entry.Username)
		fmt.Println("Password: ", entry.Password)
		fmt.Println("#==================================#")
	}
	return nil
}
