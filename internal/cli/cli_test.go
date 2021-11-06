package cli_test

import (
	"testing"

	"github.com/theprimeagen/projectizer/internal/cli"
)

func TestAddNoArgs(t *testing.T) {
    pwd := "foo/bar/baz"
    _, err := cli.New(&cli.CliArgs{
        pwd,
        []string{"add"},
    })

    if err == nil {
        t.Fatalf("Expected add to return error")
    }
}

func TestAddOneArg(t *testing.T) {
    pwd := "foo/bar/baz"
    _, err := cli.New(&cli.CliArgs{
        pwd,
        []string{"add", "arg1"},
    })

    if err == nil {
        t.Fatalf("Expected add to return error")
    }
}

func TestAddSuccess(t *testing.T) {
    pwd := "foo/bar/baz"
    config, err := cli.New(&cli.CliArgs{
        pwd,
        []string{"add", "arg1", "arg2", "arg3"},
    })

    if err != nil {
        t.Fatalf("expected add to return error")
    }

    if config.Pwd != pwd {
        t.Fatalf("expected pwd directory to be equal to %q", pwd)
    }

    args := []string{
        "arg1",
        "arg2",
        "arg3",
    }

    for i, arg := range config.AdditionalArgs {
        if arg != args[i] {
            t.Fatalf("expected arg %q to equal %q", arg, args[i])
        }
    }
}

func TestPrintWithTooManyArgs(t *testing.T) {
    pwd := "foo/bar/baz"
    _, err := cli.New(&cli.CliArgs{
        pwd,
        []string{"foo", "bar"},
    })

    if err == nil {
        t.Fatalf("expected to fail with an error")
    }
}

func TestPrintWithArg(t *testing.T) {
    pwd := "foo/bar/baz"
    config, err := cli.New(&cli.CliArgs{
        pwd,
        []string{"foo"},
    })

    if err != nil {
        t.Fatalf("expected print not error.")
    }

    if config.Cmd != "print" {
        t.Fatalf("expected cmd to be print but got %q", config.Cmd)
    }
}

func TestPrint(t *testing.T) {
    pwd := "foo/bar/baz"
    config, err := cli.New(&cli.CliArgs{
        pwd,
        []string{},
    })

    if err != nil {
        t.Fatalf("expected print not error.")
    }

    if config.Cmd != "print" {
        t.Fatalf("expected cmd to be print but got %q", config.Cmd)
    }
}
