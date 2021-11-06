package cli

import (
	"errors"
	"flag"
	"os"
)

type CliConfig struct {
	Pwd string
	Cmd string
    AdditionalArgs []string
}

type CliArgs struct {
    Pwd string;
    Args []string;
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

func New(cliArgs *CliArgs) (*CliConfig, error) {

	cmd := "print"
	args := cliArgs.Args;

	if len(args) > 0 {
		cmd = args[0]
        args = args[1:]
	}

	switch cmd {
	case "print":
		if len(args) > 0 {
			return nil, errors.New("too many arguments")
		}
	case "link":
		if len(args) == 0 {
			return nil, errors.New("please provide the path to link to")
		} else if len(args) > 1 {
			return nil, errors.New("too many arguments")
		}
	case "add":
		if len(args) < 2 {
			return nil, errors.New("please provide the command_name followed by the command for add")
		}
	default:
        return nil, errors.New("invalid command provided")
	}

	return &CliConfig{
		cliArgs.Pwd,
		cmd,
        args,
	}, nil
}
