package project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/theprimeagen/projectizer/internal/cli"
)

type Project map[string]string
type ProjectMap map[string]Project

func getProject(proj *ProjectMap, config *cli.CliArgs) {

    for k, v := range *proj {
        fmt.Printf("key[%s] value[%s]\n", k, v)
    }
}

func Run(config *cli.CliArgs) {
    // TODO: Make this better...
    home, err := os.UserHomeDir()
    if err != nil {
        log.Fatalf("%+v\n", err)
    }

    projectPath := path.Join(home, "./.project.json")
    proj, err := ioutil.ReadFile(projectPath)

    if err != nil {
        log.Fatalf("%+v\n", err)
    }

    var data ProjectMap
    err = json.Unmarshal(proj, &data)
    if err != nil {
        log.Fatalf("%+v\n", err)
    }

    getProject(&data, config)
}

