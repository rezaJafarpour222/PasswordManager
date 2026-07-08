package cli

import (
	"fmt"
	"os"
	"pass/TUI"
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

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Command ")

	dp := []TUI.DataPoint{}
	for _, entry := range a.Commands {
		dp = append(dp, TUI.DataPoint{Key: entry.Name, Value: entry.Description})
	}
	box.PrintData(dp)
}

func (a *App) Init() {
	box := TUI.NewBox(100, '╭', '╮', '╰', '╯')
	done := make(chan struct{})
	_, err := os.Stat(a.VaultPath)
	if err == nil {
		box.SetTitle("Init")
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error", Value: "Vault does exist"})
		box.PrintData(dp)
		return
	}
	_, err = os.Stat(a.MasterKeyPath)
	if err == nil {
		box.SetTitle("Init")
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error", Value: "Master key does exist"})
		box.PrintData(dp)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go Spinner(done, &wg, "Initializing the Vault and Master Key")
	err = storage.SaveMasterKey(a.MasterKeyPath)
	if err != nil {
		box.SetTitle("Init")
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error", Value: "Problem Loading Master key"})
		box.PrintData(dp)
		return
	}
	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		box.SetTitle("Init")
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error", Value: "Problem Loading Vault"})
		box.PrintData(dp)
		return
	}
	v := encryption.NewVault()
	err = storage.SaveVault(v, masterKey, a.VaultPath)
	if err != nil {
		box.SetTitle("Init")
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error", Value: "Problem Saving Vault"})
		box.PrintData(dp)
		return
	}
	close(done)
	wg.Wait()
	txt := fmt.Sprintf("%s and %s created.", a.VaultPath, a.MasterKeyPath)
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Success", Value: txt})
	box.PrintData(dp)
	fmt.Println()
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

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Vault")
	dp := []TUI.DataPoint{}
	for _, entry := range v.Entries {
		dp = append(dp, TUI.DataPoint{Key: "Service ", Value: entry.Service})
		dp = append(dp, TUI.DataPoint{Key: "Username", Value: entry.Username})
		dp = append(dp, TUI.DataPoint{Key: "Password", Value: entry.Password})
	}
	box.PrintData(dp)
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

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Entry added")
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Service ", Value: service})
	dp = append(dp, TUI.DataPoint{Key: "Username", Value: username})
	dp = append(dp, TUI.DataPoint{Key: "Password", Value: password})
	box.PrintData(dp)
	return nil
}

func (a *App) Gen(service, username string, sizeStr string) error {

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return fmt.Errorf("Password length must be a number.")
	}
	if size > 32 {
		return fmt.Errorf("Password length must be <=32.")
	}

	randomPassword, err := encryption.GenerateRandomPassword(size)
	if err != nil {
		return err
	}

	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Service ", Value: service})
	dp = append(dp, TUI.DataPoint{Key: "Username", Value: username})
	dp = append(dp, TUI.DataPoint{Key: "Password", Value: randomPassword})
	if !Confirmation("generating password", dp) {
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
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Service ", Value: service})
	if !Confirmation("Delete", dp) {
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
		box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
		box.SetTitle("Deleting " + service)
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Error", Value: err.Error()})
		box.PrintData(dp)
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

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Result For: " + service)
	dp := []TUI.DataPoint{}
	for _, entry := range entries {
		dp = append(dp, TUI.DataPoint{Key: "Service ", Value: entry.Service})
		dp = append(dp, TUI.DataPoint{Key: "Username", Value: entry.Username})
		dp = append(dp, TUI.DataPoint{Key: "Password", Value: entry.Password})
	}
	box.PrintData(dp)
	return nil
}
func (a *App) GetMasterKey() error {

	_, err := os.Stat(a.MasterKeyPath)
	if err != nil {
		return fmt.Errorf("master key file does not exist.")
	}

	masterKey, err := storage.LoadMasterKey(a.MasterKeyPath)
	if err != nil {
		return err
	}

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Key")
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Master Key", Value: masterKey})
	box.PrintData(dp)
	return nil
}

func (a *App) Export(exportPath string) error {
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "vault to: ", Value: exportPath + "/Vault.vault"})
	dp = append(dp, TUI.DataPoint{Key: "master key To: ", Value: exportPath + "/Master.key"})
	Confirmation("export", dp)
	err := storage.ExportVault(exportPath, a.VaultPath)
	if err != nil {
		return err
	}
	err = storage.ExportMasterKey(exportPath, a.MasterKeyPath)
	if err != nil {
		return err
	}
	return nil
}
