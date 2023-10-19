package main

import (
	"os"

	"github.com/bersen66/my_little_shell/pkg/shell"
)

func main() {
	config := shell.NewConfig()
	CM := make(map[string]func() shell.Command)
	CM["pwd"] = func() shell.Command {
		return &shell.Pwd{
			Config: config,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}
	CM["echo"] = func() shell.Command {
		return &shell.Echo{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}
	CM["cd"] = func() shell.Command {
		return &shell.Cd{
			Config: config,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}
	CM["ps"] = func() shell.Command {
		return &shell.Ps{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}
	CM["kill"] = func() shell.Command {
		return &shell.Kill{}
	}
	shell.Run(config, CM)
}
