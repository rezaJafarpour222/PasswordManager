package cli

import (
	"fmt"
	"strings"
)

type ParsedCommand struct {
	Name  string
	Args  []string
	Flags map[string]string
}

func (a *App) Parse(args []string) (*ParsedCommand, error) {

	if len(args) == 0 {
		return nil, fmt.Errorf("no command provided")
	}

	cmdName := args[0]

	if _, ok := a.Commands[cmdName]; !ok {
		return nil, fmt.Errorf("unknown command: %s", cmdName)
	}

	cmd := &ParsedCommand{
		Name:  cmdName,
		Args:  []string{},
		Flags: make(map[string]string),
	}

	for i := 1; i < len(args); i++ {

		arg := args[i]

		if strings.HasPrefix(arg, "-") {

			key := strings.TrimLeft(arg, "-")

			if i+1 >= len(args) {
				return nil, fmt.Errorf("missing value for flag %s", arg)
			}

			value := args[i+1]

			if strings.HasPrefix(value, "-") {
				return nil, fmt.Errorf("missing value for flag %s", arg)
			}

			cmd.Flags[key] = value
			i++
			continue
		}

		cmd.Args = append(cmd.Args, arg)
	}

	return cmd, nil
}
