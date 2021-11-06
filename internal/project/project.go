package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/theprimeagen/projectizer/internal/cli"
)

type ProjectMapJSON map[string]map[string]string
type ProjectJSON struct {
	Aliases  map[string]string `json:"aliases"`
	Projects ProjectMapJSON    `json:"projects"`
}

func getProjectPath() string {

	// TODO: Make this better...
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	return path.Join(home, "./.project.json")
}

type ProjectDataProvider interface {
	Get(config *cli.CliConfig) ([]byte, string, error)
    Set(path string, data []byte) error
}

type Project struct {
	path    string
	Project ProjectJSON
    provider ProjectDataProvider
}

type FileDataProvider struct{}

func (p *FileDataProvider) Get(config *cli.CliConfig) ([]byte, string, error) {
	projectPath := getProjectPath()
	if _, err := os.Stat(projectPath); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(projectPath, []byte("{\"aliases\": {}, \"projects\": {}}"), 0644)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
	}

    proj, err := os.ReadFile(projectPath)

    return proj, projectPath, err
}

func (p *FileDataProvider) Set(path string, data []byte) error {
    return os.WriteFile(path, data, 0644)
}

func New(config *cli.CliConfig, provider ProjectDataProvider) (*Project, error) {
	proj, projectPath, err := provider.Get(config)

	if err != nil {
        return nil, err
	}

	var data ProjectJSON
	err = json.Unmarshal(proj, &data)
	if err != nil {
        return nil, err
	}

    return &Project{path: projectPath, Project: data, provider: provider}, nil
}

func (p *Project) Save() error {
    proj, err := json.Marshal(p.Project)
    if err != nil {
        return err
    }

    p.provider.Set(p.path, proj)
	return nil
}

func (p *Project) print(projectPath string, config *cli.CliConfig) (bool, error) {
    if len(projectPath) == 0 {
        return false, errors.New("couldn't find project")
    }

    if len(config.AdditionalArgs) > 0 {
        keyName := config.AdditionalArgs[0]
        if val, ok := p.Project.Projects[projectPath][keyName]; ok {
            projJSON, _ := json.Marshal(val)
            fmt.Println(string(projJSON))
        } else {
            return false, fmt.Errorf("no key found in this project with name %q", keyName)
        }
    } else {
        projJSON, _ := json.Marshal(p.Project.Projects[projectPath])
        fmt.Println(string(projJSON))
    }

    return false, nil
}

func (p *Project) add(projectPath string, config *cli.CliConfig) (bool, error) {
    keyName := config.AdditionalArgs[0]
    value := strings.Join(config.AdditionalArgs[1:], " ")

    if projectPath == "" {
        projectPath = config.Pwd
        p.Project.Projects[projectPath] = map[string]string{
            keyName: value,
        }
    } else {
        p.Project.Projects[projectPath][keyName] = value
    }

    return true, nil
}
func (p *Project) link(projectPath string, config *cli.CliConfig) (bool, error) {
    p.Project.Aliases[config.Pwd] = config.AdditionalArgs[0]
    return true, nil
}

func (p *Project) unlink(projectPath string, config *cli.CliConfig) (bool, error) {
    _, ok := p.Project.Aliases[config.Pwd];
    if ok {
        delete(p.Project.Aliases, config.Pwd)
    }
    return true, nil
}

func (p *Project) Run(config *cli.CliConfig) (bool, error) {
	projectPath := ""
	if val, ok := p.Project.Aliases[config.Pwd]; ok {
		projectPath = val
	} else if _, ok := p.Project.Projects[config.Pwd]; ok {
		projectPath = config.Pwd
	}

	changed := false

	// run the Effing project
	switch config.Cmd {
    case "print":
        return p.print(projectPath, config)
	case "add":
        return p.add(projectPath, config)
	case "link":
        return p.link(projectPath, config)
	case "unlink":
        return p.unlink(projectPath, config)
	}

	return changed, nil
}
