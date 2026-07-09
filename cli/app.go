package cli

import (
	"pass/TUI"
)

type Command struct {
	Name        string
	Description string
	Example     string
}

type App struct {
	VaultPath     string
	MasterKeyPath string
	Version       string
	Commands      map[string]Command
}

func New(version, vaultPath, masterKeyPath string) *App {
	app := &App{
		VaultPath:     vaultPath,
		MasterKeyPath: masterKeyPath,
		Version:       version,
		Commands:      make(map[string]Command),
	}

	app.registerCommand("init", "   Create a vault", "pass init")
	app.registerCommand("add", "    Add new entry to vault", "pass add -u username -p password")
	app.registerCommand("del", "    delete entry from the vault", "pass delete service")
	app.registerCommand("get", "    Get a credential", "pass get servicename -u username")
	app.registerCommand("key", "    Get a master key", "pass key")
	app.registerCommand("list", "   List all credentials", "pass list")
	app.registerCommand("export", " Export the vault and master key ", "pass export -p /path/ShouldBe/A/Folder")
	app.registerCommand("gen", "    Generate a random password for service and username", "pass gen servicename -u username -l 12")
	app.registerCommand("version", "cli version ", "pass version")
	app.registerCommand("help", "   Help", "pass help")

	return app
}
func (a *App) Run(args []string) error {
	cmd, err := a.Parse(args[1:])
	if err != nil {
		return err
	}

	box := TUI.NewBox(60, '╭', '╮', '╰', '╯')
	switch cmd.Name {

	case "init":
		err := a.Init()
		if err != nil {
			return err
		}
	case "add":
		err := a.Add(cmd.Args[0], cmd.Flags["u"], cmd.Flags["p"])
		if err != nil {
			return err
		}
	case "get":
		err := a.GetEntry(cmd.Args[0])
		if err != nil {
			return err
		}
		return nil
	case "del":
		err := a.DeleteEntry(cmd.Args[0])
		if err != nil {
			return err
		}
		return nil

	case "gen":
		err := a.Gen(cmd.Args[0], cmd.Flags["u"], cmd.Flags["l"])
		if err != nil {
			return err
		}
		return nil
	case "list":
		err := a.List()
		if err != nil {
			return err
		}

	case "key":
		err := a.GetMasterKey()
		if err != nil {
			return err
		}
		return nil
	case "version":

		box.SetTitle("Version")
		dp := []TUI.DataPoint{}
		dp = append(dp, TUI.DataPoint{Key: "Password manager Version ", Value: a.Version})
		box.PrintData(dp)
		return nil

	case "export":
		err := a.Export(cmd.Args[0])
		if err != nil {
			return err
		}
		return nil
	case "help":
		a.Help()
		return nil
	}

	return nil
}
