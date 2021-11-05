package cli

import (
	"flag"
	"log"
	"os"
)

type CliArgs struct {
    Pwd string
    Cmd string
}

func New() *CliArgs {
    args := os.Args

    var pwd string
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("%+v\n", err)
    }

    flag.StringVar(&pwd, "pwd", cwd, "which project to get config for")
    flag.Parse();

    cmd := "print"
    if len(args) > 1 {
        cmd = args[1]
    }

    return &CliArgs{
        pwd,
        cmd,
    }
}



