package cli

import (
	"fmt"
)

type Command struct {
	Name        string
	Description string
}

type App struct {
	FilePath string
	Version  string
	Commands map[string]Command
}

func New(version string, filePath string) *App {
	app := &App{
		FilePath: filePath,
		Version:  version,
		Commands: make(map[string]Command),
	}

	app.registerCommand("init", "Create a vault")
	app.registerCommand("add", "Add new entry to vault")
	app.registerCommand("get", "Get a credential")
	app.registerCommand("list", "List all credentials")
	app.registerCommand("gen", "Generate password for a service with username")
	app.registerCommand("version", "cli version ")
	app.registerCommand("help", "Help")

	return app
}
func (a *App) Run(args []string, masterpassword string) error {

	cmd, err := a.Parse(args[1:])
	if err != nil {
		return err
	}

	switch cmd.Name {

	case "add":
		fmt.Println(cmd.Args)
		err := a.Add(cmd.Args[0], cmd.Flags["u"], cmd.Flags["p"], masterpassword)
		if err != nil {
			return err
		}

	case "get":
		fmt.Println("service:")
		return nil

	case "gen":
		fmt.Println("service:Generator")
		a.Gen(cmd.Args[0], cmd.Flags["u"], masterpassword, 32)
		return nil
	case "init":
		err := a.InitCheck()
		if err != nil {
			return err
		}
		a.Init(masterpassword)
	case "list":
		a.List(masterpassword)
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
