package cli

import (
	"fmt"
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
	app.registerCommand("del", " delete entry from the vault", "pass delete service")
	app.registerCommand("get", "    Get a credential", "pass get servicename -u username")
	app.registerCommand("list", "   List all credentials", "pass list")
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

	switch cmd.Name {

	case "add":
		fmt.Println(cmd.Args)
		err := a.Add(cmd.Args[0], cmd.Flags["u"], cmd.Flags["p"])
		if err != nil {
			return err
		}

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
	case "init":
		err := a.Init()
		if err != nil {
			return err
		}
	case "list":
		err := a.List()
		if err != nil {
			return err
		}
		return nil
	case "version":
		fmt.Println(a.Version)
		return nil
	case "help":
		a.Help()
		return nil
	}

	return nil
}
