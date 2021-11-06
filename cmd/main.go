package main

import (
	"log"

	"github.com/theprimeagen/projectizer/internal/cli"
	"github.com/theprimeagen/projectizer/internal/project"
)

func main() {
    args, err := cli.GetCLIArgs()
    if err != nil {
        log.Fatalf("%+v\n", err)
    }

    config, err := cli.New(args)
    if err != nil {
        log.Fatalf("%+v\n", err)
    }

    fileReader := project.FileDataProvider{}
    project, err := project.New(config, &fileReader)
    if err != nil {
        log.Fatalf("Project: %+v\n", err)
    }

    changed, err := project.Run(config)
    if err != nil {
        log.Fatalf("Project: %+v\n", err)
    }

    if changed {
        err = project.Save()
        if err != nil {
            log.Fatalf("%+v\n", err)
        }
    }
}

