package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bersen66/my_little_shell/pkg/utils"
)

func ConstructCommands(cm map[string]func() Command, commands []string) []Command {
	result := make([]Command, 0, len(commands))

	for _, command := range commands {
		words := strings.Split(command, " ")
		args := make([]string, 0)
		for i, _ := range words {
			words[i] = strings.Trim(words[i], " \t")
			if words[i] != "" {
				args = append(args, words[i])
			}
		}

		if builder, has := cm[args[0]]; has {
			cmd := builder()
			cmd.ParseArgs(args[1:])
			result = append(result, cmd)
		} else {
			c := &NotBuiltin{}
			c.ParseArgs(args)
			result = append(result, c)
		}

	}

	return result
}

func PipeCommands(commands []Command) {

	for i, command := range commands[:len(commands)-1] {
		out, _ := command.StdoutPipe()
		command.SetStderr(os.Stderr)
		command.Run()
		commands[i+1].SetStdin(out)
	}
	commands[len(commands)-1].SetStderr(os.Stderr)
}

func Run(config *Config, cm map[string]func() Command) {
	utils.ClearCmd()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%v ", config.UiString())
	for scanner.Scan() {
		line := scanner.Text()

		commands := ConstructCommands(cm, strings.Split(line, "|"))
		PipeCommands(commands)

		result, _ := commands[len(commands)-1].Output()
		fmt.Printf("%s", result)

		fmt.Printf("%v ", config.UiString())
	}
}
