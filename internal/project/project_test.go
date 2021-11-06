package project_test

import (
	"fmt"
	"testing"

	"github.com/theprimeagen/projectizer/internal/cli"
	"github.com/theprimeagen/projectizer/internal/project"
)

type EmptyProvider struct {
    Aliases string;
    Projects string;
}

func New() EmptyProvider {
    return EmptyProvider{
        Aliases: "{}",
        Projects: "{}",
    }
}

func (t *EmptyProvider) Get(config *cli.CliConfig) ([]byte, string, error) {
    return []byte(fmt.Sprintf("{\"aliases\": %s, \"projects\": %s}", t.Aliases, t.Projects)), "foo/bar", nil
}

func (t *EmptyProvider) Set(path string, data []byte) error {
    return nil
}

func TestAdd(t *testing.T) {
    pwd := "foo/bar/baz"
    provider := New()
    config := cli.CliConfig{
        Pwd: pwd,
        Cmd: "add",
        AdditionalArgs: []string{"foo", "bar", "baz"},
    }

    project, err := project.New(&config, &provider)

    if err != nil {
        t.Fatalf("expected new#add to not error %+v", err)
    }

    if len(project.Project.Projects) != 0 {
        t.Fatalf("expected no projects but got %d", len(project.Project.Projects))
    }

    changed, err := project.Run(&config)

    if err != nil {
        t.Fatalf("expected no error %+v", err)
    }
    if !changed {
        t.Fatalf("expected there to be a change")
    }

    if len(project.Project.Projects) == 0 {
        t.Fatalf("expected run#add to add a project.")
    }

    fmt.Printf("foo key: %+v", project.Project.Projects)
    if val, ok := project.Project.Projects[pwd]; ok {
        if val2, ok := val["foo"]; ok {
            if val2 != "bar baz" {
                t.Fatalf("expected foo to have a value of %q but got %q", "bar baz", val2)
            }
        } else {
            t.Fatalf("expected the project to have key foo")
        }
    } else {
        t.Fatalf("expected foo to be a key in the projects.")
    }
}

func TestPrintFailNoProject(t *testing.T) {
    pwd := "foo/bar/baz"
    provider := New()
    config := cli.CliConfig{
        Pwd: pwd,
        Cmd: "print",
        AdditionalArgs: []string{},
    }

    project, err := project.New(&config, &provider)

    if err != nil {
        t.Fatalf("expected new#print to not error %+v", err)
    }

    _, err = project.Run(&config)

    if err == nil {
        t.Fatalf("expected print to error since there is no project %+v", err)
    }
}

func TestPrintFailNoKey(t *testing.T) {
    pwd := "foo/bar/baz"
    provider := New()
    provider.Projects = fmt.Sprintf("{\"%s\": {}}", pwd)

    config := cli.CliConfig{
        Pwd: pwd,
        Cmd: "print",
        AdditionalArgs: []string{"foo"},
    }

    project, err := project.New(&config, &provider)

    if err != nil {
        t.Fatalf("expected new#print to not error %+v", err)
    }

    _, err = project.Run(&config)

    if err == nil {
        t.Fatalf("expected print to error since there is no project %+v", err)
    }
}
