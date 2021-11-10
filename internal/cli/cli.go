package cli

import (
	"errors"
	"flag"
	"os"
)

type CliConfig struct {
	Pwd            string
	Cmd            string
	AdditionalArgs []string
}

type CliArgs struct {
	Pwd  string
	Args []string
}

func GetCLIArgs() (*CliArgs, error) {
	var pwd string
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// project -- cwd config
	// project add <cmd_name> ... -- cwd empty config
	// project link dir ... -- cwd empty config
	flag.StringVar(&pwd, "pwd", cwd, "which project to get config for")
	flag.Parse()

	return &CliArgs{pwd, flag.Args()}, nil
}

var commands = map[string]bool{
    "unlink": true,
    "link": true,
    "add": true,
    "print": true,
    "del": true,
}

func isCommand(cmd string) bool {
    _, ok := commands[cmd]
    return ok
}

func New(cliArgs *CliArgs) (*CliConfig, error) {

	cmd := "print"
	args := cliArgs.Args

	if len(args) > 0 && isCommand(args[0]) {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "print":
		if len(args) > 1 {
			return nil, errors.New("too many arguments, print can take in 1 argument to print that key for this project")
		}
	case "link":
		if len(args) == 0 {
			return nil, errors.New("please provide the path to link to")
		} else if len(args) > 1 {
			return nil, errors.New("too many arguments")
		}
	case "unlink":
		if len(args) != 0 {
			return nil, errors.New("too many arguments")
		}
	case "del":
		if len(args) > 1 {
			return nil, errors.New("too many arguments")
		}
	case "add":
		if len(args) < 2 {
			return nil, errors.New("please provide the command_name followed by the command for add")
		}
	case "which":
		if len(args) > 0 {
			return nil, errors.New("too many additional arguments")
		}
	}

	return &CliConfig{
		cliArgs.Pwd,
		cmd,
		args,
	}, nil
}
