package main

import (
	"github.com/theprimeagen/projectizer/internal/cli"
	"github.com/theprimeagen/projectizer/internal/project"
)

func main() {
    config := cli.New()
    project.Run(config)
}

