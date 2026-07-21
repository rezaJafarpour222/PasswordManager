package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"pass/TUI"
	"pass/protocol"
	"pass/storage"
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

func (a *App) Init() error {
	done := make(chan struct{})

	req := protocol.Request{
		Action: "init",
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go TUI.Spinner(done, &wg, "Initializing the Vault and Master Key\n")

	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}

	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
	}
	close(done)
	wg.Wait()
	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Success", Value: string(res.Payload)})
	box.PrintData(dp)
	fmt.Println()
	return nil
}

func (a *App) List() error {
	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)
	go TUI.Spinner(done, &wg, "Decrypting the Vault")
	req := protocol.Request{
		Action: "list",
	}
	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}

	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
	}
	entries, err := protocol.DeserializeEntries(res)
	close(done)
	wg.Wait()
	if len(entries) == 0 {
		return fmt.Errorf("Vault is empty")
	}

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Vault")
	for _, entry := range entries {
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Service ", Value: entry.Service})
		dp = append(dp, TUI.DataPoint{Key: "Username", Value: entry.Username})
		dp = append(dp, TUI.DataPoint{Key: "Password", Value: entry.Password})
		box.PrintData(dp)
	}
	return nil
}

func (a *App) Add(service, username, password string) error {

	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)

	spinnerText := fmt.Sprintf("Adding '%s' to the vault", service)
	go TUI.Spinner(done, &wg, spinnerText)

	close(done)
	wg.Wait()
	req := protocol.Request{
		Action:   "add",
		Service:  service,
		Username: username,
		Password: password,
	}
	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}
	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
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

	req := protocol.Request{
		Action:   "gen",
		Service:  service,
		Username: username,
	}
	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}
	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
	}
	entry, err := protocol.DeserializeEntry(res)
	if err != nil {
		return err
	}

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Entry Generated")
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Service ", Value: entry.Service})
	dp = append(dp, TUI.DataPoint{Key: "Username", Value: entry.Username})
	dp = append(dp, TUI.DataPoint{Key: "Password", Value: entry.Password})
	box.PrintData(dp)
	return nil
}

func (a *App) DeleteEntry(service string) error {

	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)
	spinnerText := fmt.Sprintf("Deleting '%s' from the vault\n", service)
	dp := []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Service ", Value: service})
	if !Confirmation("Delete", dp) {
		return nil
	}
	go TUI.Spinner(done, &wg, spinnerText)

	req := protocol.Request{
		Action:  "del",
		Service: service,
	}
	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}
	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
	}
	_, err = protocol.DeserializeEntry(res)
	if err != nil {
		return err
	}
	close(done)
	wg.Wait()

	return nil
}

func (a *App) GetEntry(service string) error {

	req := protocol.Request{
		Action:  "get",
		Service: service,
	}
	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}
	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
	}
	entries, err := protocol.DeserializeEntries(res)
	if err != nil {
		return err
	}
	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Result For: " + service)
	for _, entry := range entries {
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Service ", Value: entry.Service})
		dp = append(dp, TUI.DataPoint{Key: "Username", Value: entry.Username})
		dp = append(dp, TUI.DataPoint{Key: "Password", Value: entry.Password})
		box.PrintData(dp)
	}
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
	if !Confirmation("export", dp) {
		return nil
	}
	req := protocol.Request{
		Action:     "export",
		ExportPath: exportPath,
	}
	res, err := protocol.Dial(req)
	if err != nil {
		return err
	}
	if res.Error {
		var msg string
		json.Unmarshal(res.Payload, &msg)
		return errors.New(msg)
	}

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	box.SetTitle("Export")
	dp = []TUI.DataPoint{}
	dp = append(dp, TUI.DataPoint{Key: "Result", Value: string(res.Payload)})
	box.PrintData(dp)
	return nil
}
