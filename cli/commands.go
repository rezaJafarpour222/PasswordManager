package cli

import (
	"fmt"
	"os"
	"pass/encryption"
	"pass/storage"
)

func (a *App) registerCommand(name string, desc string) {
	a.Commands[name] = Command{
		Name:        name,
		Description: desc,
	}
}
func (a *App) Help() {
	fmt.Println("Usage:")
	fmt.Println("		pass <command> [option]")
	fmt.Println()
	fmt.Println("Commands:")
	for _, cmd := range a.Commands {
		fmt.Printf("		%-12s  %s\n", cmd.Name, cmd.Description)
	}
}
func (a *App) Init(masterPassword string) {
	v := encryption.NewVault()
	err := storage.SaveVault(v, masterPassword, a.FilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Vault Created")
}

func (a *App) List(masterPassword string) error {
	v, err := storage.LoadVault(masterPassword, a.FilePath)
	if err != nil {
		panic(err)
	}
	for _, entry := range v.Entries {
		fmt.Println("Service: ", entry.Service)
		fmt.Println("Username: ", entry.Username)
		fmt.Println("Password: ", entry.Password)
		fmt.Println("#====================#")
	}
	return nil

}
func (a *App) InitCheck() error {
	_, err := os.Stat(a.FilePath)
	if err == nil {
		return fmt.Errorf("Vault does exist!")
	}
	return nil
}

func (a *App) Add(service, username, password, masterPassword string) error {
	v, err := storage.LoadVault(masterPassword, a.FilePath)
	if err != nil {
		return err
	}

	v.AddEntry(encryption.Entry{
		Service:  service,
		Username: username,
		Password: password,
	})

	return storage.SaveVault(v, masterPassword, a.FilePath)
}

func (a *App) Gen(service, username, masterPassword string, size int) error {
	randomPassword, err := encryption.GenerateRandomPassword(size)
	if err != nil {
		return err
	}
	err = a.Add(service, username, randomPassword, masterPassword)
	if err != nil {
		return err
	}
	return nil
}
