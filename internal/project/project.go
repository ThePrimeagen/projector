package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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
	path     string
	Project  ProjectJSON
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

func filePathPop(path string) string {
	// I don't like this.
	path, _ = filepath.Split(path)
	path = path[0 : len(path)-1]

	return path
}

func (p *Project) getProjectPath(base string, depth int) string {
	projectPath := ""
	for ; depth != 0; depth-- {
		if val, ok := p.Project.Aliases[base]; ok {
			projectPath = val
			break
		}

		if _, ok := p.Project.Projects[base]; ok {
			projectPath = base
			break
		}

		base = filePathPop(base)
		if len(base) == 0 {
			break
		}
	}

	return projectPath
}

func (p *Project) print(config *cli.CliConfig) (bool, error) {
	keyName := ""
	if len(config.AdditionalArgs) > 0 {
		keyName = config.AdditionalArgs[0]
	}

    toPrint := make(map[string]string)
    projectPath := p.getProjectPath(config.Pwd, -1)
	for ; len(projectPath) != 0; projectPath = p.getProjectPath(filePathPop(projectPath), -1) {
		if len(keyName) > 0 {
			if val, ok := p.Project.Projects[projectPath][keyName]; ok {
				fmt.Println(val)
				break
			}
		} else {
            for k, v := range p.Project.Projects[projectPath] {
                if _, ok := p.Project.Projects[projectPath][keyName]; !ok {
                    toPrint[k] = v
                }
            }
		}
	}

    if len(toPrint) > 0 {
        projJSON, _ := json.Marshal(toPrint)
        fmt.Println(string(projJSON))
    } else if len(projectPath) == 0 {
		return false, errors.New("couldn't find project")
	}

	return false, nil
}

func (p *Project) add(config *cli.CliConfig) (bool, error) {
	keyName := config.AdditionalArgs[0]
	value := strings.Join(config.AdditionalArgs[1:], " ")

    projectPath := p.getProjectPath(config.Pwd, 1)
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
func (p *Project) link(config *cli.CliConfig) (bool, error) {
	p.Project.Aliases[config.Pwd] = config.AdditionalArgs[0]
	return true, nil
}

func (p *Project) unlink(config *cli.CliConfig) (bool, error) {
	_, ok := p.Project.Aliases[config.Pwd]
	if ok {
		delete(p.Project.Aliases, config.Pwd)
	}
	return true, nil
}

func (p *Project) Run(config *cli.CliConfig) (bool, error) {

	// run the Effing project
	switch config.Cmd {
	case "print":
		return p.print(config)
	case "add":
		return p.add(config)
	case "link":
		return p.link(config)
	case "unlink":
		return p.unlink(config)
	default:
		return false, fmt.Errorf("file an issue, this should never happen %s", config.Cmd)
	}
}
